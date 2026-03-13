#!/bin/bash
#
# download-all-voices.sh - Download all Piper TTS voice models
#
# These voice models are used with text-to-speech.sh to generate audio
# for Remotion video compositions. For voice selection guidance, see:
#   .agents/video-producer/skills/audio-synthesis.md
#
# Usage: ./download-all-voices.sh [quality_filter]
#   quality_filter: low, medium, high, or all (default: all)
#

PIPER_DIR="${PIPER_DIR:-$HOME/piper}"
QUALITY="${1:-all}"
mkdir -p "$PIPER_DIR"

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_voice() { echo -e "${BLUE}[VOICE]${NC} $1"; }

# Voice models: lang/region/voice_name:voice_id_prefix
# voice_id_prefix is like en_US-amy, en_GB-alan, etc.
VOICES=(
    # English US
    "en/en_US/amy:en_US-amy"
    "en/en_US/arctic:en_US-arctic"
    "en/en_US/danny:en_US-danny"
    "en/en_US/hfc_female:en_US-hfc_female"
    "en/en_US/hfc_male:en_US-hfc_male"
    "en/en_US/joe:en_US-joe"
    "en/en_US/kathleen:en_US-kathleen"
    "en/en_US/kusal:en_US-kusal"
    "en/en_US/l2arctic:en_US-l2arctic"
    "en/en_US/lessac:en_US-lessac"
    "en/en_US/libritts:en_US-libritts"
    "en/en_US/ryan:en_US-ryan"
    
    # English GB
    "en/en_GB/alba:en_GB-alba"
    "en/en_GB/alan:en_GB-alan"
    "en/en_GB/cori:en_GB-cori"
    "en/en_GB/jenny_dioco:en_GB-jenny_dioco"
    "en/en_GB/semaine:en_GB-semaine"
    "en/en_GB/vctk:en_GB-vctk"
    
    # English Australia
    "en/en_AU/narayan:en_AU-narayan"
    
    # English India
    "en/en_IN/arctic:en_IN-arctic"
    
    # Spanish Spain
    "es/es_ES/carlfm:es_ES-carlfm"
    "es/es_ES/davefx:es_ES-davefx"
    "es/es_ES/sharvard:es_ES-sharvard"
    
    # Spanish Mexico
    "es/es_MX/claudia:es_MX-claudia"
    
    # French
    "fr/fr_FR/gilles:fr_FR-gilles"
    "fr/fr_FR/siwis:fr_FR-siwis"
    "fr/fr_FR/upmc:fr_FR-upmc"
    
    # German
    "de/de_DE/eva_k:de_DE-eva_k"
    "de/de_DE/thorsten:de_DE-thorsten"
    "de/de_DE/thorsten_emotional:de_DE-thorsten_emotional"
    
    # Italian
    "it/it_IT/riccardo:it_IT-riccardo"
    
    # Portuguese Brazil
    "pt/pt_BR/edresson:pt_BR-edresson"
    "pt/pt_BR/faber:pt_BR-faber"
    
    # Portuguese Portugal
    "pt/pt_PT/tugao:pt_PT-tugao"
    
    # Dutch
    "nl/nl_NL/nathalie:nl_NL-nathalie"
    "nl/nl_NL/rdh:nl_NL-rdh"
    
    # Polish
    "pl/pl_PL/darkman:pl_PL-darkman"
    "pl/pl_PL/gosia:pl_PL-gosia"
    "pl/pl_PL/mc_speech:pl_PL-mc_speech"
    
    # Russian
    "ru/ru_RU/irina:ru_RU-irina"
    "ru/ru_RU/ruslan:ru_RU-ruslan"
    
    # Chinese
    "zh/zh_CN/huayan:zh_CN-huayan"
    
    # Japanese
    "ja/ja_JP/kokoro:ja_JP-kokoro"
    
    # Korean
    "ko/ko_KO/bundletrain:ko_KO-bundletrain"
    
    # Turkish
    "tr/tr_TR/df:tr_TR-df"
    "tr/tr_TR/fahrettin:tr_TR-fahrettin"
    
    # Arabic
    "ar/ar_JO/kareem:ar_JO-kareem"
    
    # Czech
    "cs/cs_CZ/jirka:cs_CZ-jirka"
    
    # Greek
    "el/el_GR/rapunzelina:el_GR-rapunzelina"
    
    # Finnish
    "fi/fi_FI/harri:fi_FI-harri"
    
    # Hungarian
    "hu/hu_HU/anna:hu_HU-anna"
    "hu/hu_HU/berta:hu_HU-berta"
    "hu/hu_HU/imre:hu_HU-imre"
    
    # Norwegian
    "no/no_NO/talesyntese:no_NO-talesyntese"
    
    # Swedish
    "sv/sv_SE/nst:sv_SE-nst"
    
    # Ukrainian
    "uk/uk_UA/lada:uk_UA-lada"
    
    # Vietnamese
    "vi/vi_VN/25hours-single:vi_VN-25hours-single"
    "vi/vi_VN/vais1000-single:vi_VN-vais1000-single"
    
    # Icelandic
    "is/is_IS/ugla:is_IS-ugla"
    
    # Nepali
    "ne/ne_NP/google:ne_NP-google"
    
    # Serbian
    "sr/sr_RS/serbski-institute:sr_RS-serbski-institute"
    
    # Swahili
    "sw/sw_CD/lanfrica:sw_CD-lanfrica"
    
    # Georgian
    "ka/ka_GE/natia:ka_GE-natia"
    
    # Bosnian
    "bs/ba_BA/goran:bs_BA-goran"
)

# Quality levels
if [ "$QUALITY" == "all" ]; then
    QUALITIES=("low" "medium" "high")
elif [ "$QUALITY" == "high" ]; then
    QUALITIES=("high")
elif [ "$QUALITY" == "medium" ]; then
    QUALITIES=("medium" "high")
else
    QUALITIES=("low" "medium")
fi

BASE_URL="https://huggingface.co/rhasspy/piper-voices/resolve/main"

log_info "Piper Voices Downloader"
log_info "Destination: $PIPER_DIR"
log_info "Quality filter: ${QUALITIES[*]}"
log_info "Starting downloads..."
echo ""

# Create download list
> /tmp/voices_urls.txt
for voice_entry in "${VOICES[@]}"; do
    path_prefix="${voice_entry%%:*}"
    voice_id="${voice_entry##*:}"
    
    for quality in "${QUALITIES[@]}"; do
        model_file="$PIPER_DIR/${voice_id}-${quality}.onnx"
        
        # Skip if already exists
        [ -f "$model_file" ] && [ -s "$model_file" ] && continue
        
        model_url="${BASE_URL}/${path_prefix}/${quality}/${voice_id}-${quality}.onnx"
        config_url="${BASE_URL}/${path_prefix}/${quality}/${voice_id}-${quality}.onnx.json"
        
        echo "$model_url|$PIPER_DIR/${voice_id}-${quality}.onnx" >> /tmp/voices_urls.txt
        echo "$config_url|$PIPER_DIR/${voice_id}-${quality}.onnx.json" >> /tmp/voices_urls.txt
    done
done

total=$(wc -l < /tmp/voices_urls.txt)
log_info "Total files to download: $total"
echo ""

# Download
count=0
while IFS='|' read -r url output; do
    ((count++))
    voice_file=$(basename "$output")
    voice_file="${voice_file%.onnx}"
    voice_file="${voice_file%.onnx.json}"
    
    if [ -f "$output" ] && [ -s "$output" ]; then
        continue
    fi
    
    log_voice "[$count/$total] $voice_file"
    
    if curl -sLf "$url" -o "$output" 2>/dev/null; then
        if [ -f "$output" ] && [ -s "$output" ]; then
            size=$(du -h "$output" | cut -f1)
            log_info "OK: $voice_file ($size)"
        else
            rm -f "$output"
            log_info "EMPTY: $voice_file"
        fi
    else
        log_info "FAIL: $voice_file"
    fi
    
    sleep 0.3
done < /tmp/voices_urls.txt

echo ""
log_info "Download complete!"
log_info "Installed voices:"
ls "$PIPER_DIR"/*.onnx 2>/dev/null | wc -l
echo ""
log_info "Available voices:"
ls "$PIPER_DIR"/*.onnx 2>/dev/null | xargs -n1 basename 2>/dev/null | sed 's/.onnx//' | sort
