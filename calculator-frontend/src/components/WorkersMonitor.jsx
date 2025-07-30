// src/components/WorkersMonitor.jsx
import { useState, useEffect } from 'react';

const WORKER_PORTS = [8081, 8082, 8083];

export default function WorkersMonitor() {
  const [workers, setWorkers] = useState(
    WORKER_PORTS.map(port => ({
      port,
      status: '⏳', // 🟡 Пока неизвестно
      loading: true
    }))
  );

  // Проверить статус одного воркера
  const checkWorker = async (port) => {
    try {
      // ✅ Меняем на POST запрос
      const res = await fetch(`http://localhost:${port}/health`, {
        method: 'POST', // ← Было GET, стало POST
        mode: 'cors',
        cache: 'no-cache',
        headers: {
          'Content-Type': 'application/json'
        }
      });

      return res.ok ? '✅' : '❌';
    } catch (err) {
      return '❌';
    }
  };

  // Опросить всех воркеров
  const pollWorkers = async () => {
    const results = await Promise.all(
      WORKER_PORTS.map(async (port) => {
        const status = await checkWorker(port);
        return { port, status, loading: false };
      })
    );

    setWorkers(results);
  };

  // Первый опрос и повтор каждые 3 секунды
  useEffect(() => {
    pollWorkers();
    const interval = setInterval(pollWorkers, 3000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="container">
      <div className="card">
        <h2 style={{ color: '#bb86fc' }}>📡 Мониторинг воркеров</h2>
        <p style={{ color: '#aaa', marginBottom: '20px' }}>
          Автоматически обновляется каждые 3 секунды
        </p>

        <div style={{
          display: 'flex',
          flexDirection: 'column',
          gap: '15px'
        }}>
          {workers.map((worker) => (
            <div
              key={worker.port}
              style={{
                display: 'flex',
                alignItems: 'center',
                padding: '15px',
                background: '#1e1e24',
                borderRadius: '10px',
                border: '1px solid #333'
              }}
            >
              <div style={{
                fontSize: '1.5rem',
                width: '40px',
                textAlign: 'center'
              }}>
                {worker.status}
              </div>
              <div style={{ flex: 1 }}>
                <div style={{ color: '#e0e0e0', fontWeight: 'bold' }}>
                  Worker на порту <code style={{ color: '#03dac6' }}>{worker.port}</code>
                </div>
                <div style={{ color: '#aaa', fontSize: '0.9rem' }}>
                  {worker.status === '✅'
                    ? '🟢 Активен'
                    : worker.status === '❌'
                      ? '🔴 Недоступен'
                      : '🟡 Проверка...'}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}