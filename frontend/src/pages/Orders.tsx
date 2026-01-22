import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

const API = import.meta.env.VITE_API_URL;

type Order = {
  ID: number;
  Product: string;
  Status: string;
};

type OrderWithAsset = {
  order: Order;
  asset: {
    ID: number;
    OrderID: number;
    LogoURL: string;
    MockupURL: string;
    AIGenerated: boolean;
    AIPrompt: string;
  };
};

function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, { bg: string, text: string }> = {
    CREATED: { bg: "#fef3c7", text: "#92400e" },
    MOCKUP_GENERATED: { bg: "#dbeafe", text: "#1e40af" },
    APPROVED: { bg: "#dcfce7", text: "#166534" },
    READY_FOR_FULFILLMENT: { bg: "#e9d5ff", text: "#7c3aed" },
  };

  const colorScheme = colors[status] || { bg: "#f3f4f6", text: "#374151" };

  return (
    <span
      style={{
        padding: "4px 10px",
        borderRadius: 12,
        fontSize: 12,
        color: colorScheme.text,
        backgroundColor: colorScheme.bg,
        fontWeight: 500,
      }}
    >
      {status.replaceAll("_", " ")}
    </span>
  );
}


export default function Orders() {
  const [ordersWithAssets, setOrdersWithAssets] = useState<OrderWithAsset[]>([]);

  useEffect(() => {
    // First get all orders
    fetch(`${API}/orders`)
      .then(res => res.json())
      .then((ordersList: Order[]) => {
        // Then fetch each order with its asset details
        Promise.all(
          ordersList.map(order => 
            fetch(`${API}/orders/${order.ID}`)
              .then(res => res.json())
              .catch(() => ({ order, asset: null }))
          )
        ).then(setOrdersWithAssets);
      });
  }, []);

  return (
    <div style={{ maxWidth: 1000, margin: "0 auto" }}>
      <h2 style={{ marginBottom: 24, color: "#111827" }}>Orders</h2>

      <div style={{ 
        backgroundColor: "white", 
        borderRadius: 12, 
        overflow: "hidden",
        border: "1px solid #e5e7eb"
      }}>
        <table style={{ 
          width: "100%", 
          borderCollapse: "collapse",
          fontSize: 14
        }}>
          <thead>
            <tr style={{ backgroundColor: "#f9fafb" }}>
              <th style={{ padding: "12px 16px", textAlign: "left", fontWeight: 600 }}>ID</th>
              <th style={{ padding: "12px 16px", textAlign: "left", fontWeight: 600 }}>Product</th>
              <th style={{ padding: "12px 16px", textAlign: "left", fontWeight: 600 }}>Status</th>
              <th style={{ padding: "12px 16px", textAlign: "left", fontWeight: 600 }}>Mockup</th>
              <th style={{ padding: "12px 16px", textAlign: "left", fontWeight: 600 }}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {ordersWithAssets.map(({ order: o, asset }) => (
              <tr key={o.ID} style={{ borderTop: "1px solid #f3f4f6" }}>
                <td style={{ padding: "12px 16px", fontWeight: 500 }}>#{o.ID}</td>
                <td style={{ padding: "12px 16px" }}>{o.Product}</td>
                <td style={{ padding: "12px 16px" }}>
                  <StatusBadge status={o.Status} />
                </td>
                <td style={{ padding: "12px 16px" }}>
                  {o.Status === "CREATED" || !asset?.MockupURL ? (
                    <span style={{ 
                      color: "#6b7280", 
                      fontSize: 12,
                      fontStyle: "italic" 
                    }}>
                      {o.Status === "CREATED" ? "Pending generation" : "No mockup"}
                    </span>
                  ) : (
                    asset.MockupURL.endsWith('.html') ? (
                      <div style={{
                        width: 50,
                        height: 50,
                        backgroundColor: "#eff6ff",
                        border: "1px solid #dbeafe",
                        borderRadius: 6,
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "center",
                        fontSize: 20
                      }}>
                        🤖
                      </div>
                    ) : (
                      <img
                        src={`${API}${asset.MockupURL}`}
                        alt="mockup"
                        style={{ 
                          width: 50, 
                          height: 50, 
                          borderRadius: 6, 
                          objectFit: "cover",
                          border: "1px solid #e5e7eb"
                        }}
                        onError={(e) => {
                          e.currentTarget.style.display = "none";
                        }}
                      />
                    )
                  )}
                </td>
                <td style={{ padding: "12px 16px" }}>
                  <Link 
                    to={`/orders/${o.ID}`}
                    style={{
                      padding: "6px 12px",
                      backgroundColor: "#3b82f6",
                      color: "white",
                      textDecoration: "none",
                      borderRadius: 6,
                      fontSize: 12,
                      fontWeight: 500
                    }}
                  >
                    View Details
                  </Link>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
