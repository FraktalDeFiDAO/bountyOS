#!/bin/bash
#
# text-to-speech.sh - Convert text files to speech using Piper TTS
#
# Output WAV files are designed for use in Remotion compositions.
# Place generated audio in your Remotion project's public/audio/ directory
# and reference via staticFile('audio/yourfile.wav') in your composition.
# For ElevenLabs (higher quality), see: .agents/video-producer/skills/audio-synthesis.md
#
# Usage: ./text-to-speech.sh [input_file] [output_dir] [voice]
#
# Examples:
#   ./text-to-speech.sh media/yt/transcript.txt media/remotion/public/audio
#   ./text-to-speech.sh media/yt/transcript.txt media/yt/audio en_US-lessac-medium
#

set -e

# Default configuration
PIPER_DIR="${PIPER_DIR:-$HOME/piper}"
PIPER_MODEL="${PIPER_MODEL:-en_US-lessac-medium}"
OUTPUT_DIR="${2:-./audio}"
VOICE="${3:-$PIPER_MODEL}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if piper is available
check_piper() {
    if command -v piper-tts &> /dev/null; then
        PIPER_CMD="piper-tts"
        log_info "Found piper-tts in PATH"
    elif command -v piper &> /dev/null; then
        PIPER_CMD="piper"
        log_info "Found piper in PATH"
    elif [ -f "$PIPER_DIR/piper-tts" ]; then
        PIPER_CMD="$PIPER_DIR/piper-tts"
        log_info "Found piper-tts at $PIPER_DIR"
    elif [ -f "$PIPER_DIR/piper" ]; then
        PIPER_CMD="$PIPER_DIR/piper"
        log_info "Found piper at $PIPER_DIR"
    else
        log_error "Piper TTS not found. Please install piper-tts or set PIPER_DIR"
        log_info "Install from: https://github.com/rhasspy/piper"
        exit 1
    fi
}

# Check if voice model exists
check_voice_model() {
    local model_file="$PIPER_DIR/${VOICE}.onnx"
    local config_file="$PIPER_DIR/${VOICE}.onnx.json"
    
    # Try without directory prefix if voice name includes it
    if [ ! -f "$model_file" ]; then
        # Extract just the voice name (e.g., en_US-lessac-medium from full path)
        local voice_name=$(basename "$VOICE")
        model_file="$PIPER_DIR/${voice_name}.onnx"
        config_file="$PIPER_DIR/${voice_name}.onnx.json"
    fi
    
    if [ ! -f "$model_file" ]; then
        log_error "Voice model not found: $model_file"
        log_info "Available models should be in: $PIPER_DIR"
        log_info "Download models from: https://huggingface.co/rhasspy/piper-voices"
        exit 1
    fi
    
    if [ ! -f "$config_file" ]; then
        log_warn "Config file not found: $config_file (may still work)"
    fi
    
    PIPER_MODEL_FILE="$model_file"
}

# Process a single text file
process_file() {
    local input_file="$1"
    local output_file="$2"
    
    if [ ! -f "$input_file" ]; then
        log_error "Input file not found: $input_file"
        return 1
    fi
    
    log_info "Processing: $input_file"
    log_info "Output: $output_file"
    
    # Run piper TTS
    "$PIPER_CMD" \
        --model "$PIPER_MODEL_FILE" \
        --output_file "$output_file" \
        < "$input_file"
    
    if [ $? -eq 0 ]; then
        log_info "Successfully created: $output_file"
    else
        log_error "Failed to process: $input_file"
        return 1
    fi
}

# Process text with sentence-by-sentence breakdown (better quality for long files)
process_file_chunked() {
    local input_file="$1"
    local output_file="$2"
    local chunk_size="${3:-500}"  # lines per chunk
    
    if [ ! -f "$input_file" ]; then
        log_error "Input file not found: $input_file"
        return 1
    fi
    
    log_info "Processing (chunked): $input_file"
    
    local temp_dir=$(mktemp -d)
    local temp_wav="$temp_dir/output.wav"
    local line_count=0
    local chunk_num=0
    local wav_files=()
    
    # Split file into chunks and process each
    while IFS= read -r line || [ -n "$line" ]; do
        # Skip empty lines
        [ -z "$line" ] && continue
        
        echo "$line" >> "$temp_dir/chunk_$chunk_num.txt"
        ((line_count++))
        
        if [ $line_count -ge $chunk_size ]; then
            wav_files+=("$temp_dir/chunk_$chunk_num.wav")
            "$PIPER_CMD" \
                --model "$PIPER_MODEL_FILE" \
                --output_file "$temp_dir/chunk_$chunk_num.wav" \
                < "$temp_dir/chunk_$chunk_num.txt"
            
            ((chunk_num++))
            line_count=0
            rm -f "$temp_dir/chunk_$chunk_num.txt"
        fi
    done < "$input_file"
    
    # Process remaining lines
    if [ $line_count -gt 0 ]; then
        wav_files+=("$temp_dir/chunk_$chunk_num.wav")
        "$PIPER_CMD" \
            --model "$PIPER_MODEL_FILE" \
            --output_file "$temp_dir/chunk_$chunk_num.wav" \
            < "$temp_dir/chunk_$chunk_num.txt"
    fi
    
    # Concatenate all WAV files
    if [ ${#wav_files[@]} -gt 1 ]; then
        log_info "Concatenating ${#wav_files[@]} audio chunks..."
        
        # Create file list for ffmpeg
        local file_list="$temp_dir/files.txt"
        for wav in "${wav_files[@]}"; do
            echo "file '$wav'" >> "$file_list"
        done
        
        # Use ffmpeg to concatenate
        if command -v ffmpeg &> /dev/null; then
            ffmpeg -y -f concat -safe 0 -i "$file_list" -c copy "$output_file" 2>/dev/null
        else
            # Fallback: use sox if available
            if command -v sox &> /dev/null; then
                sox "${wav_files[@]}" "$output_file"
            else
                log_warn "ffmpeg/sox not found. Only first chunk saved."
                cp "${wav_files[0]}" "$output_file"
            fi
        fi
    elif [ ${#wav_files[@]} -eq 1 ]; then
        cp "${wav_files[0]}" "$output_file"
    fi
    
    # Cleanup
    rm -rf "$temp_dir"
    
    log_info "Successfully created: $output_file"
}

# Process all text files in a directory
process_directory() {
    local input_dir="$1"
    
    log_info "Processing all .txt files in: $input_dir"
    
    for file in "$input_dir"/*.txt; do
        [ -f "$file" ] || continue
        
        local basename=$(basename "$file" .txt)
        local output_file="$OUTPUT_DIR/${basename}.wav"
        
        process_file "$file" "$output_file"
    done
}

# Show available voices
list_voices() {
    log_info "Looking for voice models in: $PIPER_DIR"
    
    if [ -d "$PIPER_DIR" ]; then
        find "$PIPER_DIR" -name "*.onnx" -type f 2>/dev/null | while read model; do
            basename "$model" .onnx
        done
    else
        log_warn "Piper directory not found: $PIPER_DIR"
    fi
}

# Show usage
show_usage() {
    cat << EOF
Piper TTS Text-to-Speech Converter

Usage: $0 [OPTIONS] [INPUT]

Arguments:
  INPUT           Input text file or directory (default: process all in current dir)
  OUTPUT_DIR      Output directory for audio files (default: ./audio)
  VOICE           Voice model name (default: en_US-lessac-medium)

Options:
  -h, --help      Show this help message
  -l, --list      List available voice models
  -c, --chunked   Process in chunks (better for long files)
  -s, --size N    Chunk size in lines (default: 500)

Environment Variables:
  PIPER_DIR       Path to piper installation (default: \$HOME/piper)
  PIPER_MODEL     Default voice model (default: en_US-lessac-medium)

Examples:
  $0 transcript.txt ./audio
  $0 transcript.txt ./audio en_US-lessac-medium
  $0 -c -s 300 long-book.txt ./audiobook
  $0 ./texts/ ./audio en_GB-alan-medium

Download voices from: https://huggingface.co/rhasspy/piper-voices
EOF
}

# Main execution
main() {
    local input=""
    local output_dir=""
    local voice=""
    local chunked=false
    local chunk_size=500

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_usage
                exit 0
                ;;
            -l|--list)
                check_piper
                list_voices
                exit 0
                ;;
            -c|--chunked)
                chunked=true
                shift
                ;;
            -s|--size)
                chunk_size="$2"
                shift 2
                ;;
            -*)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
            *)
                if [ -z "$input" ]; then
                    input="$1"
                elif [ -z "$output_dir" ]; then
                    output_dir="$1"
                else
                    voice="$1"
                fi
                shift
                ;;
        esac
    done

    # Set defaults
    OUTPUT_DIR="${output_dir:-./audio}"
    VOICE="${voice:-$PIPER_MODEL}"

    # Check piper installation
    check_piper
    check_voice_model

    # Create output directory
    mkdir -p "$OUTPUT_DIR"

    # Process input
    if [ -z "$input" ]; then
        log_error "No input file specified"
        show_usage
        exit 1
    elif [ -f "$input" ]; then
        local output_file="$OUTPUT_DIR/$(basename "$input" .txt).wav"

        if [ "$chunked" = true ]; then
            process_file_chunked "$input" "$output_file" "$chunk_size"
        else
            process_file "$input" "$output_file"
        fi
    elif [ -d "$input" ]; then
        process_directory "$input"
    else
        log_error "Input not found: $input"
        exit 1
    fi
}

main "$@"
