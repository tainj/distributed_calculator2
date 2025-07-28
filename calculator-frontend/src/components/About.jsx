// src/components/About.jsx
export default function About() {
  return (
    <div className="container">
      <div className="card">
        <h2 style={{ 
          color: '#bb86fc',
          marginBottom: '20px'
        }}>
          🧮 О проекте
        </h2>
        
        <p style={{ color: '#e0e0e0', fontSize: '1.1rem' }}>
          <strong>CalcPro</strong> — это распределённая система на <strong style={{ color: '#00ADD8' }}>Go</strong>, использующая <strong>Apache Kafka</strong>, <strong>Redis</strong>, <strong>PostgreSQL</strong> и <strong>gRPC</strong>.
        </p>

        <h3 style={{ color: '#03dac6', marginTop: '30px' }}>🚀 Как это работает</h3>
        <div style={{
          background: '#1e1e24',
          padding: '20px',
          borderRadius: '10px',
          margin: '15px 0'
        }}>
          <pre style={{
            color: '#e0e0e0',
            fontSize: '0.9rem',
            overflow: 'auto',
            margin: '0'
          }}>
{`[Frontend] → [Gateway] → [Kafka]
                    ↓
         [Worker 1] [Worker 2] [Worker N]
                    ↓
               [Redis Cache]
                    ↓
            [PostgreSQL History]`}
          </pre>
        </div>

        <h3 style={{ color: '#03dac6', marginTop: '30px' }}>🔧 Технологии</h3>
        <ul style={{ 
          marginLeft: '20px', 
          lineHeight: '1.8',
          color: '#e0e0e0'
        }}>
          <li><strong style={{ color: '#00ADD8' }}>Go</strong> — бэкенд, воркеры, gRPC</li>
          <li><strong>Apache Kafka</strong> — асинхронная очередь задач</li>
          <li><strong style={{ color: '#DC382D' }}>Redis</strong> — кэширование промежуточных результатов</li>
          <li><strong style={{ color: '#336791' }}>PostgreSQL</strong> — история вычислений</li>
          <li><strong style={{ color: '#61DAFB' }}>React</strong> — современный фронтенд</li>
          <li><strong>JWT</strong> — безопасная авторизация</li>
        </ul>

        <h3 style={{ color: '#03dac6', marginTop: '30px' }}>🧠 Особенности</h3>
        <ul style={{ 
          marginLeft: '20px', 
          lineHeight: '1.8',
          color: '#e0e0e0'
        }}>
          <li>Унарный минус через <code style={{ color: '#bb86fc' }}>~</code>: <code>~(~5) = 5</code></li>
          <li>Безопасный парсинг — без <code style={{ color: '#cf6679' }}>eval()</code></li>
          <li>Асинхронная обработка через очередь</li>
          <li>Масштабируемость — добавляй воркеров одной командой</li>
          <li>Деление на ноль → ошибка, не бесконечность</li>
        </ul>

        <div style={{ 
          textAlign: 'right', 
          marginTop: '40px',
          paddingTop: '20px',
          borderTop: '1px solid #2d2d3a'
        }}>
          <p style={{ 
            margin: '0', 
            color: '#aaa',
            fontStyle: 'italic'
          }}>
            Сделано с ❤️ к Go<br/>
            <span style={{ color: '#bb86fc' }}>@tainj</span>
          </p>
        </div>
      </div>
    </div>
  );
}