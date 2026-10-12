"""Create original MIT sample media; requires Python 3 and FFmpeg for GIF/MP4."""
from pathlib import Path
import math
import struct
import subprocess
import wave

root = Path(__file__).resolve().parents[1]
assets = root / 'assets'
gallery = root / 'internal/gallery/assets'
(assets / 'sample.svg').write_text('''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 360" role="img" aria-label="Geometric sample">
<rect width="640" height="360" fill="#eef2ff"/><circle cx="470" cy="110" r="55" fill="#335cff"/>
<path d="M0 360 220 100 440 360Z" fill="#14213a"/><path d="M220 360 440 160 640 360Z" fill="#8ca4ff"/>
</svg>\n''')
(assets / 'sample-before.svg').write_text((assets / 'sample.svg').read_text().replace('#335cff', '#96523e').replace('#8ca4ff', '#e7a38b'))
with wave.open(str(gallery / 'audit-tone.wav'), 'wb') as audio:
    audio.setparams((1, 2, 8000, 16000, 'NONE', 'not compressed'))
    audio.writeframes(b''.join(struct.pack('<h', int(1600 * math.sin(2 * math.pi * 440 * i / 8000))) for i in range(16000)))
subprocess.run(['ffmpeg', '-y', '-loglevel', 'error', '-f', 'lavfi', '-i',
                'testsrc2=size=160x90:rate=10:duration=2', '-threads', '1',
                '-loop', '0', str(assets / 'sample.gif')], check=True)
subprocess.run(['ffmpeg', '-y', '-loglevel', 'error', '-f', 'lavfi', '-i',
                'testsrc2=size=160x90:rate=10:duration=2', '-c:v', 'libx264',
                '-threads', '1', '-pix_fmt', 'yuv420p', '-movflags', '+faststart',
                str(gallery / 'audit-video.mp4')], check=True)
assert (assets / 'sample.gif').stat().st_size > 0
assert (gallery / 'audit-video.mp4').stat().st_size > 0
