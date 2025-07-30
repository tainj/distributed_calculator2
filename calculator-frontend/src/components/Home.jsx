// src/components/Home.jsx
export default function Home() {
  const token = localStorage.getItem('token');

  return (
    <div className="container">
      <div className="card" style={{ textAlign: 'center', padding: '40px 30px' }}>
        <h1 style={{ color: '#bb86fc', fontSize: '2.5rem', marginBottom: '10px' }}>
          🚀 CalcFlow
        </h1>
        <p style={{ fontSize: '1.1rem', color: '#aaa', marginBottom: '30px' }}>
          Распределённый калькулятор на Go
        </p>

        <div style={{
          display: 'flex',
          justifyContent: 'center',
          gap: '20px',
          margin: '30px 0',
          flexWrap: 'wrap'
        }}>
          <Stat value="Go" label="Backend" />
          <Stat value="React" label="Frontend" />
          <Stat value="Kafka" label="Очередь" />
          <Stat value="JWT" label="Авторизация" />
        </div>

        <p style={{ margin: '20px 0', color: '#bbb', fontStyle: 'italic' }}>
          "Вычисления — это искусство. А мы — художники." 💫
        </p>

        {!token && (
          <p style={{ color: '#cf6679', marginTop: '20px' }}>
            Чтобы начать — <a href="/login" style={{ color: '#bb86fc' }}>войдите</a>
          </p>
        )}

        {token && (
          <a href="/calc" className="btn-primary" style={{ marginTop: '20px' }}>
            Перейти к калькулятору
          </a>
        )}
      </div>
    </div>
  );
}

function Stat({ value, label }) {
  return (
    <div style={{
      background: '#2d2d3a',
      padding: '15px 20px',
      borderRadius: '10px',
      minWidth: '100px'
    }}>
      <div style={{ fontSize: '1.2rem', fontWeight: 'bold', color: '#bb86fc' }}>{value}</div>
      <div style={{ fontSize: '0.9rem', color: '#aaa' }}>{label}</div>
    </div>
  );
}