# DEPRECATED: This MoviePy script is superseded by the Remotion project.
# Use media/remotion/ for all new video production.
# Remotion provides superior animation, audio sync, and React-based composition.
# Run: cd media/remotion && npx remotion render src/index.ts RidimaVideo out/video.mp4
#
# This file is kept for reference only. Do not use for new content.

from moviepy import TextClip, ColorClip, concatenate_videoclips, CompositeVideoClip
from moviepy.video.fx import FadeIn
import os

# Create directory if it doesn't exist
os.makedirs('media/ridima_bounty/video', exist_ok=True)

# Configuration
width, height = 1080, 1920  # Portrait (TikTok/X/Instagram Reel)
fps = 30
font = '/usr/share/fonts/TTF/DejaVuSans.ttf'  # Using absolute path for system font

# Brand Colors (Ridima Blue & White)
RIDIMA_BLUE = (0, 43, 154)   # #002B9A
SCAM_RED = (200, 0, 0)       # Deep red for "scam" warnings
SUCCESS_GREEN = (0, 150, 0)  # Professional green
WHITE = (255, 255, 255)

def create_scene(text, duration, bg_color=RIDIMA_BLUE, font_size=80, text_color=WHITE, transition_time=0.5):
    bg = ColorClip(size=(width, height), color=bg_color, duration=duration)
    
    # Text with padding and shadow/outline effect simulation via method='caption'
    txt = TextClip(
        text=text, 
        font_size=font_size, 
        color=text_color, 
        font=font, 
        size=(int(width * 0.85), None), 
        method='caption',
        text_align='center'
    )
    txt = txt.with_duration(duration).with_position('center')
    
    scene = CompositeVideoClip([bg, txt])
    return scene.with_effects([FadeIn(duration=transition_time)])

# Scenes based on the refined script
scenes = [
    create_scene("Tired of P2P scams?", 3, bg_color=SCAM_RED, font_size=100),
    create_scene("Complex crypto stressing you out?", 3, bg_color=SCAM_RED, font_size=90),
    create_scene("Introducing\nRidima OTC", 4, bg_color=RIDIMA_BLUE, font_size=130),
    create_scene("The safest way for Nigerians to trade.", 3.5, bg_color=RIDIMA_BLUE, font_size=80),
    create_scene("Instant SOL to Naira.", 3.5, bg_color=RIDIMA_BLUE, font_size=110),
    create_scene("Premium rates.\nNo P2P.", 3, bg_color=SUCCESS_GREEN, font_size=100),
    create_scene("No scams. No stress.", 3, bg_color=RIDIMA_BLUE, font_size=100),
    create_scene("24/7 VIP support.", 3, bg_color=RIDIMA_BLUE, font_size=90),
    create_scene("Digital assets to cash,\nanytime!", 5, bg_color=RIDIMA_BLUE, font_size=120),
    create_scene("Join the Ridima community today!\n#RidimaOTC", 4, bg_color=RIDIMA_BLUE, font_size=70)
]

# Concatenate with crossfades (padding creates overlap)
final_video = concatenate_videoclips(scenes, padding=-0.5, method="compose")

# Write final file
output_path = 'media/ridima_bounty/video/ridima_animation_v2.mp4'
final_video.write_videofile(output_path, fps=fps, codec='libx264', audio=False)

print(f"Video rendered successfully: {output_path}")
