// src/components/Examples.jsx
import { useState, useEffect } from 'react';
import api from '../services/api';

export default function Examples() {
  const [examples, setExamples] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchExamples = async () => {
      try {
        const res = await api.post('/v1/examples', {});
        setExamples(res.data.examples || []);
      } catch (err) {
        console.error('Ошибка загрузки примеров', err);
      } finally {
        setLoading(false);
      }
    };
    fetchExamples();
  }, []);

  // Форматируем дату: 27.07.2025, 20:18
  const formatDate = (isoString) => {
    try {
      const date = new Date(isoString);
      return date.toLocaleString('ru-RU', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch {
      return 'Неизвестно';
    }
  };

  if (loading) {
    return (
      <div className="container">
        <div className="card">
          <p>Загрузка истории...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container">
      <div className="card">
        <h2 style={{ color: '#bb86fc' }}>История вычислений</h2>

        {examples.length === 0 ? (
          <p>Нет записей</p>
        ) : (
          <ul style={{ listStyle: 'none', padding: 0 }}>
            {examples.map((ex) => (
              <li
                key={ex.id}
                style={{
                  padding: '15px',
                  margin: '10px 0',
                  background: '#2d2d3a',
                  borderRadius: '8px',
                  borderLeft: '4px solid #9c27b0',
                  position: 'relative'
                }}
              >
                {/* Иконка статуса */}
                <div style={{
                  position: 'absolute',
                  top: '15px',
                  right: '15px',
                  fontSize: '1.2rem'
                }}>
                  {ex.calculated === false ? (
                    '⏳' // В процессе
                  ) : ex.result !== undefined ? (
                    '✅' // Успешно
                  ) : (
                    '❌' // Ошибка
                  )}
                </div>

                {/* Выражение */}
                <div style={{ marginBottom: '8px' }}>
                  <code style={{
                    background: '#1e1e24',
                    padding: '4px 8px',
                    borderRadius: '4px',
                    fontSize: '16px'
                  }}>
                    {ex.expression}
                  </code>
                </div>

                {/* Результат или ошибка */}
                <div style={{ fontSize: '14px', color: '#aaa', marginBottom: '8px' }}>
                  {ex.calculated === false ? (
                    <em>Ожидает вычисления...</em>
                  ) : ex.result !== undefined ? (
                    <strong style={{ color: '#bb86fc' }}>Результат: {ex.result}</strong>
                  ) : (
                    <span style={{ color: '#cf6679' }}>Ошибка: деление на ноль</span>
                  )}
                </div>

                {/* Дата */}
                <small style={{ color: '#777' }}>
                  {formatDate(ex.createdAt)}
                </small>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}