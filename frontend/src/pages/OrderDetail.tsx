import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

const API = import.meta.env.VITE_API_URL;

type Order = {
  ID: number;
  Product: string;
  Status: string;
};

export default function OrderDetail() {
  const { ID } = useParams<{ ID: string }>();
  const [order, setOrder] = useState<Order | null>(null);

  const [label, setLabel] = useState({
      name: "",
      address: "",
      city: "",
      state: "",
      zip: ""
    });


  const load = () => {
    fetch(`${API}/orders/${ID}`)
      .then(res => res.json())
      .then(setOrder);
  };

  useEffect(load, [ID]);

  const approve = async () => {
    await fetch(`${API}/orders/${ID}/approve`, { method: "POST" });
    load();
  };

  const generateLabel = async () => {
  const res = await fetch(`${API}/orders/${ID}/label`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(label)
  });

  const data = await res.json();
  window.open(`${API}${data.label}`, "_blank");
};


  if (!order) return <p>Loading...</p>;

  return (
    <div style={{ maxWidth: 600 }}>
      <h2>Order #{order.ID}</h2>
      <p>Status: {order.Status}</p>

      <img
        src={`${API}/mockups/order_${order.ID}.png`}
        width={300}
        style={{ borderRadius: 12, marginTop: 16 }}
        alt="Mockup"
      />

      {order.Status === "MOCKUP_GENERATED" && (
        <button
          onClick={approve}
          style={{
            display: "block",
            marginTop: 16,
            padding: "10px 16px",
            background: "#3B82F6",
            color: "white",
            border: "none",
            borderRadius: 8
          }}
        >
          Approve for Fulfillment
        </button>
      )}

          <input placeholder="Name" onChange={e => setLabel({...label, name:e.target.value})}/>
        <input placeholder="Address" onChange={e => setLabel({...label, address:e.target.value})}/>
        <input placeholder="City" onChange={e => setLabel({...label, city:e.target.value})}/>
        <input placeholder="State" onChange={e => setLabel({...label, state:e.target.value})}/>
        <input placeholder="ZIP" onChange={e => setLabel({...label, zip:e.target.value})}/>


      <button
        onClick={generateLabel}
        style={{
          display: "block",
          marginTop: 16,
          padding: "10px 16px",
          background: "#10B981",
          color: "white",
          border: "none",
          borderRadius: 8
        }}
      >
        Generate Shipping Label
      </button>
    </div>
  );
}
