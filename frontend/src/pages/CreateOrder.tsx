import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

const API = import.meta.env.VITE_API_URL;

export default function CreateOrder() {
  const navigate = useNavigate();
  const [product, setProduct] = useState("Hoodie");
  const [color, setColor] = useState("black");
  const [size, setSize] = useState("M");
  const [logo, setLogo] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);
  const [useAI, setUseAI] = useState(false);
  const [aiPrompt, setAIPrompt] = useState("");
  const [availableColors, setAvailableColors] = useState<string[]>([]);

  useEffect(() => {
    // Fetch available colors
    fetch(`${API}/colors`)
      .then(res => res.json())
      .then(data => setAvailableColors(data.colors || []))
      .catch(err => console.error('Failed to load colors:', err));
  }, []);

  const submit = async () => {
  if (!useAI && !logo) return alert("Upload a logo or use AI generation");
  if (useAI && !aiPrompt.trim()) return alert("Enter an AI prompt for mockup generation");

  setLoading(true);

  let logoUrl = "";

  // 1. Upload logo if not using AI
  if (!useAI && logo) {
    const logoForm = new FormData();
    logoForm.append("file", logo);

    const uploadRes = await fetch(`${API}/upload/logo`, {
      method: "POST",
      body: logoForm,
    });

    if (!uploadRes.ok) {
      setLoading(false);
      alert("Logo upload failed");
      return;
    }

    const { url } = await uploadRes.json();
    logoUrl = url;
  }

  // 2. Create order
  const orderRes = await fetch(`${API}/orders`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      product,
      color,
      size,
      logoUrl: logoUrl,
      useAI: useAI,
      aiPrompt: aiPrompt,
    }),
  });

  setLoading(false);

  if (!orderRes.ok) {
    alert("Order creation failed");
    return;
  }

  const orderData = await orderRes.json();
  
  // Redirect to order detail page
  navigate(`/orders/${orderData.order.ID}`);
};


  return (
    <div style={{ maxWidth: 600, margin: "0 auto" }}>
      <h2 style={{ marginBottom: 24, color: "#111827" }}>Create Custom Order</h2>

      <div style={{ marginBottom: 20 }}>
        <label style={{ display: "block", marginBottom: 8, fontWeight: 500 }}>Product</label>
        <select 
          value={product} 
          onChange={e => setProduct(e.target.value)}
          style={{ 
            width: "100%", 
            padding: "8px 12px", 
            border: "1px solid #d1d5db", 
            borderRadius: 6,
            fontSize: 14
          }}
        >
          <option>Hoodie</option>
          <option>T-Shirt</option>
        </select>
      </div>

      <div style={{ marginBottom: 20 }}>
        <label style={{ display: "block", marginBottom: 8, fontWeight: 500 }}>Color</label>
        <select 
          value={color} 
          onChange={e => setColor(e.target.value)}
          style={{ 
            width: "100%", 
            padding: "8px 12px", 
            border: "1px solid #d1d5db", 
            borderRadius: 6,
            fontSize: 14
          }}
        >
          {availableColors.map(colorOption => (
            <option key={colorOption} value={colorOption}>
              {colorOption.charAt(0).toUpperCase() + colorOption.slice(1)}
            </option>
          ))}
        </select>
      </div>

      <div style={{ marginBottom: 20 }}>
        <label style={{ display: "block", marginBottom: 8, fontWeight: 500 }}>Size</label>
        <select 
          value={size} 
          onChange={e => setSize(e.target.value)}
          style={{ 
            width: "100%", 
            padding: "8px 12px", 
            border: "1px solid #d1d5db", 
            borderRadius: 6,
            fontSize: 14
          }}
        >
          <option value="XS">XS</option>
          <option value="S">S</option>
          <option value="M">M</option>
          <option value="L">L</option>
          <option value="XL">XL</option>
          <option value="XXL">XXL</option>
        </select>
      </div>

      {/* AI Toggle */}
      <div style={{ 
        marginBottom: 20, 
        padding: 16, 
        backgroundColor: "#f3f4f6", 
        borderRadius: 8,
        border: "1px solid #e5e7eb"
      }}>
        <label style={{ 
          display: "flex", 
          alignItems: "center", 
          cursor: "pointer",
          fontWeight: 500,
          marginBottom: 12
        }}>
          <input
            type="checkbox"
            checked={useAI}
            onChange={e => setUseAI(e.target.checked)}
            style={{ marginRight: 8 }}
          />
          🤖 Use FREE AI to generate mockup
        </label>
        
        {useAI ? (
          <div>
            <label style={{ display: "block", marginBottom: 8, fontSize: 14, color: "#6b7280" }}>
              Describe your design vision
            </label>
            <textarea
              value={aiPrompt}
              onChange={e => setAIPrompt(e.target.value)}
              placeholder="e.g., 'with a minimalist geometric logo on the chest', 'featuring a vintage band logo', 'with a nature-inspired design'"
              rows={3}
              style={{ 
                width: "100%", 
                padding: "8px 12px", 
                border: "1px solid #d1d5db", 
                borderRadius: 6,
                fontSize: 14,
                resize: "vertical"
              }}
            />
            <div style={{ fontSize: 12, color: "#6b7280", marginTop: 4 }}>
              💡 Tip: Be specific about the style, placement, and theme you want<br/>
              🆓 Powered by FREE Hugging Face Stable Diffusion - no API key needed!
            </div>
          </div>
        ) : (
          <div>
            <label style={{ display: "block", marginBottom: 8, fontSize: 14, color: "#6b7280" }}>
              Upload Logo
            </label>
            <input 
              type="file" 
              onChange={e => setLogo(e.target.files?.[0] || null)}
              accept="image/*"
              style={{ 
                width: "100%", 
                padding: "8px 12px", 
                border: "1px solid #d1d5db", 
                borderRadius: 6,
                fontSize: 14
              }}
            />
          </div>
        )}
      </div>

      <button 
        onClick={submit} 
        disabled={loading}
        style={{
          width: "100%",
          padding: "12px 24px",
          backgroundColor: loading ? "#9ca3af" : "#3b82f6",
          color: "white",
          border: "none",
          borderRadius: 6,
          fontSize: 16,
          fontWeight: 500,
          cursor: loading ? "not-allowed" : "pointer",
          transition: "background-color 0.2s"
        }}
      >
        {loading ? "Creating..." : "Create Order"}
      </button>
    </div>
  );
}
