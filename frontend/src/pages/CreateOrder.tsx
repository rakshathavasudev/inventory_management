import { useState } from "react";

const API = import.meta.env.VITE_API_URL;

export default function CreateOrder() {
  const [product, setProduct] = useState("Hoodie");
  const [color, setColor] = useState("");
  const [size, setSize] = useState("");
  const [logo, setLogo] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);

  const submit = async () => {
  if (!logo) return alert("Upload a logo");

  setLoading(true);

  // 1. Upload logo
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
      logoUrl: url,
    }),
  });

  setLoading(false);

  if (!orderRes.ok) {
    alert("Order creation failed");
    return;
  }

  alert("Order created successfully");
};


  return (
    <div>
      <h2>Create Custom Order</h2>

      <div>
        <label>Product</label><br />
        <select value={product} onChange={e => setProduct(e.target.value)}>
          <option>Hoodie</option>
          <option>T-Shirt</option>
        </select>
      </div>

      <div>
        <label>Color</label><br />
        <input value={color} onChange={e => setColor(e.target.value)} />
      </div>

      <div>
        <label>Size</label><br />
        <input value={size} onChange={e => setSize(e.target.value)} />
      </div>

      <div>
        <label>Logo</label><br />
        <input type="file" onChange={e => setLogo(e.target.files?.[0] || null)} />
      </div>

      <button onClick={submit} disabled={loading}>
        {loading ? "Creating..." : "Create Order"}
      </button>
    </div>
  );
}
