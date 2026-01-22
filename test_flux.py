import torch
from diffusers import DiffusionPipeline
import os

print("Testing FLUX.2-klein-9B model...")
print(f"PyTorch version: {torch.__version__}")
print(f"CUDA available: {torch.cuda.is_available()}")

if torch.cuda.is_available():
    print(f"CUDA device: {torch.cuda.get_device_name()}")
    device = "cuda"
else:
    print("Using CPU (this will be slow)")
    device = "cpu"

try:
    print("Loading FLUX.2-klein-9B pipeline...")
    # Use CPU if CUDA not available, and float32 for CPU compatibility
    if device == "cuda":
        pipe = DiffusionPipeline.from_pretrained("black-forest-labs/FLUX.2-klein-9B", torch_dtype=torch.bfloat16, device_map="cuda")
    else:
        pipe = DiffusionPipeline.from_pretrained("black-forest-labs/FLUX.2-klein-9B", torch_dtype=torch.float32)
        pipe = pipe.to(device)
    
    print("Pipeline loaded successfully!")
    
    prompt = "Blue hoodie with geometric patterns, professional product photography, clean white background, studio lighting, high quality, detailed, realistic, e-commerce style, front view, centered composition, photorealistic, 4k resolution"
    print(f"Generating image with prompt: {prompt}")
    
    # Generate image
    image = pipe(prompt).images[0]
    
    # Save the image
    output_path = "test_flux2_output.png"
    image.save(output_path)
    print(f"Image saved to: {output_path}")
    
except Exception as e:
    print(f"Error: {e}")
    print("Make sure you have the required packages installed:")
    print("pip install torch torchvision diffusers transformers accelerate")
    print("\nIf you get authentication errors, you may need to:")
    print("1. Create a Hugging Face account")
    print("2. Accept the model license at: https://huggingface.co/black-forest-labs/FLUX.2-klein-9B")
    print("3. Login with: huggingface-cli login")