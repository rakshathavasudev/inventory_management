import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

const API = import.meta.env.VITE_API_URL;

type Order = {
  ID: number;
  Product: string;
  Color: string;
  Size: string;
  Status: string;
};

type Asset = {
  ID: number;
  OrderID: number;
  LogoURL: string;
  MockupURL: string;
  AIGenerated: boolean;
  AIPrompt: string;
};

type OrderResponse = {
  order: Order;
  asset: Asset;
};

export default function OrderDetail() {
  const { ID } = useParams<{ ID: string }>();
  const [orderData, setOrderData] = useState<OrderResponse | null>(null);
  const [mockupLoading, setMockupLoading] = useState(false);
  const [labelLoading, setLabelLoading] = useState(false);

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
      .then(setOrderData);
  };

  useEffect(load, [ID]);

  const generateMockup = async () => {
    if (!orderData) return;
    
    setMockupLoading(true);
    
    try {
      const response = await fetch(`${API}/orders/${ID}/mockup`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          logoUrl: orderData.asset?.LogoURL || "",
          aiPrompt: orderData.asset?.AIPrompt || "",
          useAI: orderData.asset?.AIGenerated || false,
        }),
      });

      if (response.ok) {
        // Reload the order data to show the new mockup
        load();
      } else {
        const errorData = await response.json();
        alert(`Failed to generate mockup: ${errorData.error || 'Unknown error'}`);
      }
    } catch (error) {
      alert("Error generating mockup");
    } finally {
      setMockupLoading(false);
    }
  };

  const approve = async () => {
    await fetch(`${API}/orders/${ID}/approve`, { method: "POST" });
    load();
  };

  const generateLabel = async () => {
    setLabelLoading(true);
    
    try {
      const res = await fetch(`${API}/orders/${ID}/label`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(label)
      });

      if (res.ok) {
        const data = await res.json();
        window.open(`${API}${data.label}`, "_blank");
      } else {
        const errorData = await res.json();
        alert(`Failed to generate label: ${errorData.error || 'Unknown error'}`);
      }
    } catch (error) {
      alert("Error generating label");
    } finally {
      setLabelLoading(false);
    }
  };


  if (!orderData) return <p>Loading...</p>;

  const { order } = orderData;
  // asset might be null/empty for orders without assets yet

  return (
    <div style={{ maxWidth: 600, margin: "0 auto", padding: "0 16px" }}>
      <h2>Order #{order.ID}</h2>
      
      <div style={{ 
        backgroundColor: "#f9fafb", 
        padding: 16, 
        borderRadius: 8, 
        marginBottom: 20,
        border: "1px solid #e5e7eb"
      }}>
        <div style={{ marginBottom: 8 }}>
          <strong>Product:</strong> {order.Product}
        </div>
        <div style={{ marginBottom: 8 }}>
          <strong>Color:</strong> {order.Color}
        </div>
        <div style={{ marginBottom: 8 }}>
          <strong>Size:</strong> {order.Size}
        </div>
        <div style={{ marginBottom: 8 }}>
          <strong>Status:</strong> <span style={{ 
            padding: "4px 8px", 
            borderRadius: 4, 
            backgroundColor: order.Status === "APPROVED" ? "#dcfce7" : "#fef3c7",
            color: order.Status === "APPROVED" ? "#166534" : "#92400e",
            fontSize: 12,
            fontWeight: 500
          }}>
            {order.Status}
          </span>
        </div>
        
        {orderData.asset?.AIGenerated && (
          <div style={{ 
            marginTop: 12, 
            padding: 12, 
            backgroundColor: "#eff6ff", 
            borderRadius: 6,
            border: "1px solid #dbeafe"
          }}>
            <div style={{ 
              display: "flex", 
              alignItems: "center", 
              marginBottom: 8,
              color: "#1e40af",
              fontWeight: 500
            }}>
              🤖 AI-Generated Mockup
            </div>
            <div style={{ fontSize: 14, color: "#374151" }}>
              <strong>Prompt:</strong> "{orderData.asset.AIPrompt}"
            </div>
          </div>
        )}
      </div>

      {orderData.asset?.MockupURL ? (
        <div style={{ marginBottom: 20 }}>
          {orderData.asset.MockupURL.endsWith('.html') ? (
            // Display HTML mockups in an iframe
            <div style={{ 
              border: "1px solid #e5e7eb", 
              borderRadius: 12, 
              overflow: "hidden",
              backgroundColor: "#f9fafb"
            }}>
              <iframe
                src={`${API}${orderData.asset.MockupURL}`}
                width="100%"
                height="600"
                style={{ border: "none", borderRadius: 12 }}
                title="AI Mockup Preview"
              />
            </div>
          ) : (
            // Display regular image mockups
            <img
              src={`${API}${orderData.asset.MockupURL}`}
              width={300}
              style={{ borderRadius: 12, border: "1px solid #e5e7eb" }}
              alt="Product Mockup"
            />
          )}
          <div style={{ marginTop: 12 }}>
            <button
              onClick={generateMockup}
              disabled={mockupLoading}
              style={{
                padding: "8px 16px",
                background: "#6b7280",
                color: "white",
                border: "none",
                borderRadius: 6,
                fontSize: 14,
                cursor: mockupLoading ? "not-allowed" : "pointer",
                opacity: mockupLoading ? 0.6 : 1
              }}
            >
              {mockupLoading ? "Regenerating..." : "🔄 Regenerate Mockup"}
            </button>
            {orderData.asset.MockupURL.endsWith('.html') && (
              <div style={{ 
                marginTop: 8, 
                fontSize: 12, 
                color: "#6b7280",
                fontStyle: "italic"
              }}>
                💡 This is a preview mockup. Follow the setup instructions above to enable real AI image generation.
              </div>
            )}
          </div>
        </div>
      ) : (
        <div style={{ 
          marginBottom: 20,
          padding: 20,
          backgroundColor: "#f3f4f6",
          borderRadius: 12,
          border: "2px dashed #d1d5db",
          textAlign: "center"
        }}>
          <div style={{ fontSize: 48, marginBottom: 12 }}>🎨</div>
          <h3 style={{ margin: "0 0 12px 0", color: "#374151" }}>Ready to Generate Your Mockup!</h3>
          <p style={{ margin: "0 0 16px 0", color: "#6b7280" }}>
            {orderData.asset?.AIGenerated 
              ? `AI Prompt: "${orderData.asset.AIPrompt}"`
              : "Click below to generate a mockup for this order"
            }
          </p>
          <button
            onClick={generateMockup}
            disabled={mockupLoading}
            style={{
              padding: "12px 24px",
              background: mockupLoading ? "#9ca3af" : "#3b82f6",
              color: "white",
              border: "none",
              borderRadius: 8,
              fontSize: 16,
              fontWeight: 500,
              cursor: mockupLoading ? "not-allowed" : "pointer",
              transition: "background-color 0.2s"
            }}
          >
            {mockupLoading ? (
              orderData.asset?.AIGenerated ? "🤖 AI Generating..." : "⚙️ Generating..."
            ) : (
              orderData.asset?.AIGenerated ? "🤖 Generate AI Mockup" : "⚙️ Generate Mockup"
            )}
          </button>
          {orderData.asset?.AIGenerated && (
            <div style={{ 
              marginTop: 12, 
              fontSize: 12, 
              color: "#6b7280" 
            }}>
              🆓 FREE AI generation • First request may take 20-30 seconds
            </div>
          )}
        </div>
      )}

      {order.Status === "MOCKUP_GENERATED" && (
        <button
          onClick={approve}
          style={{
            display: "block",
            marginBottom: 20,
            padding: "12px 20px",
            background: "#3B82F6",
            color: "white",
            border: "none",
            borderRadius: 8,
            fontWeight: 500,
            cursor: "pointer"
          }}
        >
          Approve for Fulfillment
        </button>
      )}

      <div style={{ 
        backgroundColor: "#f9fafb", 
        padding: 16, 
        borderRadius: 8,
        border: "1px solid #e5e7eb",
        maxWidth: "100%",
        boxSizing: "border-box"
      }}>
        <h3 style={{ marginBottom: 16, color: "#111827", margin: "0 0 16px 0" }}>Shipping Information</h3>
        
        <div style={{ 
          marginBottom: 12, 
          padding: 8, 
          backgroundColor: "#eff6ff", 
          borderRadius: 6, 
          fontSize: 12, 
          color: "#1e40af",
          lineHeight: 1.4
        }}>
          💡 Tip: Leave fields empty to use default values (Customer, 123 Main Street, Anytown, CA 12345)
        </div>
        
        <div style={{ 
          display: "grid", 
          gap: 12,
          width: "100%",
          overflow: "hidden"
        }}>
          <input 
            placeholder="Full Name (default: Customer)" 
            value={label.name}
            onChange={e => setLabel({...label, name: e.target.value})}
            style={{ 
              width: "100%",
              padding: "8px 12px", 
              border: "1px solid #d1d5db", 
              borderRadius: 6,
              fontSize: 14,
              boxSizing: "border-box"
            }}
          />
          <input 
            placeholder="Street Address (default: 123 Main Street)" 
            value={label.address}
            onChange={e => setLabel({...label, address: e.target.value})}
            style={{ 
              width: "100%",
              padding: "8px 12px", 
              border: "1px solid #d1d5db", 
              borderRadius: 6,
              fontSize: 14,
              boxSizing: "border-box"
            }}
          />
          <div style={{ 
            display: "grid", 
            gridTemplateColumns: "repeat(auto-fit, minmax(120px, 1fr))",
            gap: 8,
            width: "100%"
          }}>
            <input 
              placeholder="City (default: Anytown)" 
              value={label.city}
              onChange={e => setLabel({...label, city: e.target.value})}
              style={{ 
                width: "100%",
                padding: "8px 12px", 
                border: "1px solid #d1d5db", 
                borderRadius: 6,
                fontSize: 14,
                boxSizing: "border-box",
                minWidth: 0
              }}
            />
            <input 
              placeholder="State (default: CA)" 
              value={label.state}
              onChange={e => setLabel({...label, state: e.target.value})}
              style={{ 
                width: "100%",
                padding: "8px 12px", 
                border: "1px solid #d1d5db", 
                borderRadius: 6,
                fontSize: 14,
                boxSizing: "border-box",
                minWidth: 0
              }}
            />
            <input 
              placeholder="ZIP (default: 12345)" 
              value={label.zip}
              onChange={e => setLabel({...label, zip: e.target.value})}
              style={{ 
                width: "100%",
                padding: "8px 12px", 
                border: "1px solid #d1d5db", 
                borderRadius: 6,
                fontSize: 14,
                boxSizing: "border-box",
                minWidth: 0
              }}
            />
          </div>
        </div>

        <button
          onClick={generateLabel}
          disabled={labelLoading}
          style={{
            marginTop: 16,
            padding: "12px 20px",
            background: labelLoading ? "#9ca3af" : "#10B981",
            color: "white",
            border: "none",
            borderRadius: 8,
            fontWeight: 500,
            cursor: labelLoading ? "not-allowed" : "pointer",
            opacity: labelLoading ? 0.6 : 1,
            transition: "all 0.2s"
          }}
        >
          {labelLoading ? "📄 Generating Label..." : "📄 Generate Shipping Label"}
        </button>
      </div>
    </div>
  );
}
