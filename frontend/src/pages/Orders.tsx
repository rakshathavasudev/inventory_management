import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

const API = import.meta.env.VITE_API_URL;

type Order = {
  ID: number;
  Product: string;
  Status: string;
};

function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    CREATED: "#9CA3AF",
    MOCKUP_GENERATED: "#3B82F6",
    APPROVED: "#10B981",
    READY_FOR_FULFILLMENT: "#8B5CF6",
  };

  return (
    <span
      style={{
        padding: "4px 10px",
        borderRadius: 12,
        fontSize: 12,
        color: "white",
        backgroundColor: colors[status] || "#6B7280",
      }}
    >
      {status.replaceAll("_", " ")}
    </span>
  );
}


export default function Orders() {
  const [orders, setOrders] = useState<Order[]>([]);

  useEffect(() => {
    fetch(`${API}/orders`)
      .then(res => res.json())
      .then(setOrders);
  }, []);

  return (
    <div>
      <h2>Orders</h2>

      <table border={1} cellPadding={8}>
        <thead>
          <tr>
            
            <th>ID</th>
            <th>Product</th>
            <th>Status</th>
            <th></th>
            <th>Mockup</th>
          </tr>
        </thead>
        <tbody>
          {orders.map(o => (
            <tr key={o.ID}>
              <td>{o.ID}</td>
              <td>{o.Product}</td>
              <td><StatusBadge status={o.Status} /></td>
              <td>
                <Link to={`/orders/${o.ID}`}>View</Link>
              </td>
              <td>
              <img
                src={`${API}/mockups/order_${o.ID}.png`}
                alt="mockup"
                style={{ width: 60, borderRadius: 6 }}
                onError={(e) => (e.currentTarget.style.display = "none")}
              />
            </td>

            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
