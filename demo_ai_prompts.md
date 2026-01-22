# FREE AI Mockup Generation Demo Prompts

Here are some example prompts you can try with the **FREE** AI mockup generation feature powered by Hugging Face Stable Diffusion:

## 🆓 Completely Free!

This uses Hugging Face's free inference API - no payment required! You can optionally get a free API token for higher rate limits, but it works without any setup.

## Minimalist Designs
- "clean minimalist logo with geometric shapes on the chest"
- "simple black and white typography design"
- "modern abstract symbol in the center"

## Nature & Outdoor
- "mountain silhouette with pine trees across the chest"
- "ocean wave pattern in blue tones"
- "forest landscape with hiking trail"

## Vintage & Retro
- "vintage band logo with distressed text effect"
- "retro 80s neon style design"
- "classic americana eagle with stars"

## Tech & Gaming
- "pixel art style gaming logo"
- "circuit board pattern design"
- "futuristic holographic effect logo"

## Art & Creative
- "watercolor splash design in rainbow colors"
- "hand-drawn sketch style illustration"
- "graffiti-style street art design"

## Business & Professional
- "corporate logo with clean typography"
- "professional emblem with shield design"
- "elegant monogram with decorative elements"

## Tips for Better Results

1. **Be Specific**: Include details about placement, style, and colors
2. **Mention the Product**: The system automatically includes product info, but you can be more specific
3. **Style Keywords**: Use terms like "minimalist", "vintage", "modern", "hand-drawn"
4. **Placement**: Specify "on the chest", "across the front", "small corner logo"
5. **Colors**: Mention color preferences or "monochrome", "colorful", "earth tones"

## Example Complete Prompts

### For a Black Hoodie
"minimalist white geometric logo on the chest with clean lines and modern typography"

### For a Red T-Shirt  
"vintage band logo in distressed white text with rock and roll aesthetic"

### For a Navy Hoodie
"mountain landscape silhouette in white with pine trees and stars"

## What the AI Generates

The FREE Hugging Face Stable Diffusion model creates professional product mockups that show:
- The specified design on the chosen product
- Realistic lighting and shadows
- Clean background suitable for e-commerce
- Professional photography style
- Proper scale and placement

## Rate Limits & Performance

- **Without API Key**: ~1000 requests/month (plenty for testing!)
- **With Free API Key**: Much higher limits
- **Model Loading**: First request may take 20-30 seconds (model cold start)
- **Subsequent Requests**: Usually 3-10 seconds

## Getting a Free API Key (Optional)

1. Go to https://huggingface.co/settings/tokens
2. Create a free account
3. Generate a new token (read access is enough)
4. Add it to your `.env` file as `HUGGINGFACE_API_KEY=your_token_here`

## Fallback Behavior

If AI generation fails (API issues, rate limits, etc.), the system automatically falls back to a creative HTML mockup that visualizes your prompt.