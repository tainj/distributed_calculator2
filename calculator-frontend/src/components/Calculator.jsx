// src/components/Calculator.jsx
import { useState } from 'react';
import api from '../services/api';

export default function Calculator() {
  const [expr, setExpr] = useState('');
  const [result, setResult] = useState('');
  const [loading, setLoading] = useState(false);
  const [showInfo, setShowInfo] = useState(false); // ✅ Для модального окна
  const token = localStorage.getItem('token');

  // Добавить символ
  const append = (value) => {
    setExpr(prev => prev + value);
  };

  // Очистить
  const clear = () => {
    setExpr('');
    setResult('');
  };

  // Отправить на сервер
  const calculate = async () => {
    if (!expr.trim()) return;
    if (!token) return alert('Войдите, чтобы вычислять!');

    setLoading(true);
    setResult('📤 Отправка выражения...'); // ✅ Показываем сразу

    try {
      const res = await api.post('/v1/calculate', { expression: expr });
      const taskId = res.data.taskId;

      setResult('⏳ Задача отправлена, ждём результат...'); // ✅ Показываем после получения taskId

      let attempts = 0;
      const max = 15;
      const interval = setInterval(async () => {
        attempts++;
        try {
          const res2 = await api.post('/v1/result', { task_id: taskId });
          if (res2.data.hasOwnProperty('value')) {
            setResult(`= ${res2.data.value}`);
            clearInterval(interval);
            setLoading(false);
          } else if (res2.data.hasOwnProperty('error')) {
            if (!res2.data.error.includes('not found')) {
              setResult(`❌ ${res2.data.error}`);
              clearInterval(interval);
              setLoading(false);
            } else {
              setResult(`⏳ Задача ещё в обработке... (${attempts}/${max})`); // ✅ Промежуточный статус
            }
          }
        } catch (err) {
          if (attempts >= max) {
            setResult('⚠️ Ошибка соединения');
            clearInterval(interval);
            setLoading(false);
          } else {
            setResult(`🔁 Попытка ${attempts}/${max}...`);
          }
        }
        if (attempts >= max) {
          setResult('⏰ Время вышло');
          clearInterval(interval);
          setLoading(false);
        }
      }, 1000);
    } catch (err) {
      setResult('❌ ' + (err.response?.data?.error || 'Ошибка'));
      setLoading(false);
    }
  };

  // Кнопки
  const buttons = [
    ['C', '(', ')', 'info'],
    ['^', '/', '*', '-'],
    ['7', '8', '9', '+'],
    ['4', '5', '6', '~'],
    ['1', '2', '3', '0'],
    ['=']
  ];

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'flex-start', // или 'center' если хочешь по центру
      alignItems: 'center',
      paddingTop: '40px', // Отступ сверху, можно менять
      paddingBottom: '40px',
      boxSizing: 'border-box'
    }}>
      <div className="card" style={{ 
        maxWidth: '400px', 
        width: '100%',
        margin: 0
      }}>
        <h2 style={{ color: '#bb86fc', textAlign: 'center' }}>🧮 Калькулятор</h2>

        {!token ? (
          <p style={{ color: '#cf6679', textAlign: 'center' }}>
            Чтобы вычислять — <a href="/login" style={{ color: '#bb86fc' }}>войдите</a>
          </p>
        ) : (
          <>
            {/* Поле ввода */}
            <input
              className="input-field"
              type="text"
              value={expr}
              onChange={(e) => setExpr(e.target.value)}
              placeholder="Введите выражение"
              style={{ fontSize: '1.4rem', textAlign: 'right', height: '60px' }}
            />

            {/* Результат */}
            {result && (
              <div style={{
                marginTop: '10px',
                padding: '12px',
                background: '#2d2d3a',
                borderRadius: '8px',
                fontSize: '1.3rem',
                textAlign: 'right',
                color: result.startsWith('❌') || result.startsWith('⚠️') || result.startsWith('⏰') 
                  ? '#cf6679' : '#4caf50',
                minHeight: '24px'
              }}>
                {result}
              </div>
            )}

            {/* Клавиатура */}
            <div style={{ marginTop: '20px' }}>
              {buttons.map((row, i) => (
                <div key={i} style={{
                  display: 'flex',
                  gap: '10px',
                  marginBottom: '10px'
                }}>
                  {row.map((btn) => (
                    <button
                      key={btn}
                      onClick={() => {
                        if (btn === '=') calculate();
                        else if (btn === 'C') clear();
                        else if (btn === 'info') setShowInfo(true);
                        else append(btn);
                      }}
                      className={
                        btn === '=' ? 'btn-primary' :
                        btn === 'C' ? 'btn-clear' :
                        btn === 'info' ? 'btn-info' :
                        'btn-calc'
                      }
                      disabled={loading && btn === '='}
                      style={{
                        ...(btn === '=' 
                          ? {
                              flex: '1',
                              padding: '20px 0',
                              fontSize: '1.4rem',
                              borderRadius: '8px',
                              border: 'none',
                              cursor: 'pointer',
                              background: 'linear-gradient(135deg, #9c27b0, #673ab7)',
                              color: 'white',
                              fontWeight: 'bold'
                            }
                          : {
                              flex: '1',
                              padding: '20px 0',
                              fontSize: '1.4rem',
                              borderRadius: '8px',
                              border: 'none',
                              cursor: 'pointer',
                              background: btn === 'C'
                                ? '#cf6679'
                                : btn === 'info'
                                  ? '#03dac6'
                                  : '#2d2d3a',
                              color: 'white',
                              fontWeight: 'bold',
                              aspectRatio: '1/1',
                              display: 'flex',
                              alignItems: 'center',
                              justifyContent: 'center'
                            }
                        )
                      }}
                    >
                      {btn}
                    </button>
                  ))}
                </div>
              ))}
            </div>
          </>
        )}
      </div>

      {/* Модальное окно для info */}
      {showInfo && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          background: 'rgba(0, 0, 0, 0.7)',
          display: 'flex',
          alignItems: 'flex-start',
          justifyContent: 'center',
          zIndex: 1000,
          paddingTop: '80px'
        }} onClick={() => setShowInfo(false)}>
          <div style={{
            background: '#1a1a24',
            padding: '25px',
            borderRadius: '12px',
            maxWidth: '400px',
            width: '90%',
            border: '1px solid #333',
            boxShadow: '0 10px 30px rgba(0, 0, 0, 0.5)'
          }} onClick={(e) => e.stopPropagation()}>
            <h3 style={{ color: '#bb86fc', marginBottom: '15px' }}>ℹ️ Поддерживаемые операции</h3>
            <ul style={{ 
              color: '#e0e0e0', 
              lineHeight: '1.8',
              paddingLeft: '20px'
            }}>
              <li><code>+</code>, <code>-</code>, <code>*</code>, <code>/</code> — базовые</li>
              <li><code>^</code> — возведение в степень</li>
              <li><code>~</code> — унарный минус: <code>~5</code>, <code>~(~9)</code></li>
              <li>Скобки: <code>( )</code></li>
            </ul>
            <p style={{ 
              color: '#aaa', 
              fontSize: '0.9rem',
              marginTop: '20px',
              fontStyle: 'italic'
            }}>
              Пример: <code style={{ color: '#03dac6' }}>~(2 + 3) * 4^2</code>
            </p>
            <button
              onClick={() => setShowInfo(false)}
              className="btn-primary"
              style={{
                width: '100%',
                padding: '12px',
                marginTop: '20px',
                fontSize: '1rem'
              }}
            >
              Закрыть
            </button>
          </div>
        </div>
      )}

      {/* Стили */}
      <style jsx>{`
        .btn-calc {
          background: #1e1e24;
          color: #e0e0e0;
          border: 1px solid #444;
        }
        .btn-calc:hover {
          background: #2d2d3a;
        }
        .btn-clear {
          background: #cf6679;
        }
        .btn-clear:hover {
          background: #b00020;
        }
        .btn-info {
          background: #03dac6;
        }
        .btn-info:hover {
          background: #018786;
        }
      `}</style>
    </div>
  );
}