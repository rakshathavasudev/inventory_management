import { Routes, Route, Link, useLocation } from "react-router-dom";
import CreateOrder from "./pages/CreateOrder";
import Orders from "./pages/Orders";
import OrderDetail from "./pages/OrderDetail";

function NavLink({
  to,
  label,
}: {
  to: string;
  label: string;
}) {
  const location = useLocation();
  const active = location.pathname === to;

  return (
    <Link
      to={to}
      style={{
        marginRight: 16,
        textDecoration: "none",
        fontWeight: active ? 600 : 500,
        color: active ? "#111827" : "#6b7280",
      }}
    >
      {label}
    </Link>
  );
}

export default function App() {
  return (
    <div
      style={{
        minHeight: "100vh",
        background: "#f8fafc",
        fontFamily: "Inter, system-ui, sans-serif",
      }}
    >
      {/* Header */}
      <header
        style={{
          background: "white",
          borderBottom: "1px solid #e5e7eb",
          padding: "16px 24px",
        }}
      >
        <div
          style={{
            maxWidth: 1200,
            margin: "0 auto",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
          }}
        >
          <div style={{ fontWeight: 700 }}>Shipify</div>

          <nav>
            <NavLink to="/" label="Create Order" />
            <NavLink to="/orders" label="Orders" />
          </nav>
        </div>
      </header>

      {/* Main Content */}
      <main
        style={{
          maxWidth: 1200,
          margin: "0 auto",
          padding: "24px",
        }}
      >
        <Routes>
          <Route path="/" element={<CreateOrder />} />
          <Route path="/orders" element={<Orders />} />
          <Route path="/orders/:ID" element={<OrderDetail />} />
        </Routes>
      </main>
    </div>
  );
}
