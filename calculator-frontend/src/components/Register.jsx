// src/components/Register.jsx
import { useState } from 'react';
import api from '../services/api';
import { Link, useNavigate } from 'react-router-dom';

export default function Register() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false); // для глазика
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const res = await api.post('/v1/register', { email, password });

      if (res.data.success) {
        alert('Регистрация успешна! Войдите.');
        navigate('/login');
      } else {
        setError(res.data.error || 'Ошибка регистрации');
      }
    } catch (err) {
      setError('Не удалось подключиться к серверу');
      console.error('Register error:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      backgroundColor: '#0f0f14',
      padding: '20px'
    }}>
      <div className="card" style={{ width: '100%', maxWidth: '400px' }}>
        <h2 style={{ color: '#bb86fc', textAlign: 'center', marginBottom: '20px' }}>
          Регистрация
        </h2>

        {error && (
          <p style={{
            backgroundColor: '#3a1b1b',
            color: '#cf6679',
            padding: '12px',
            borderRadius: '8px',
            marginBottom: '15px',
            fontSize: '14px',
            textAlign: 'center'
          }}>
            {error}
          </p>
        )}

        <form onSubmit={handleSubmit}>
          {/* Поле Email */}
          <div style={{ position: 'relative', marginBottom: '15px' }}>
            <span style={{
              position: 'absolute',
              left: '15px',
              top: '50%',
              transform: 'translateY(-50%)',
              color: '#aaa'
            }}>
              ✉️
            </span>
            <input
              className="input-field"
              type="email"
              placeholder="Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              style={{ paddingLeft: '45px' }}
            />
          </div>

          {/* Поле Пароль с глазиком */}
          <div style={{ position: 'relative', marginBottom: '20px' }}>
            <span style={{
              position: 'absolute',
              left: '15px',
              top: '50%',
              transform: 'translateY(-50%)',
              color: '#aaa'
            }}>
              🔐
            </span>
            <input
              className="input-field"
              type={showPassword ? 'text' : 'password'}
              placeholder="Пароль"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              style={{ paddingLeft: '45px', paddingRight: '70px' }}
            />
            {/* Кнопка "глазик" */}
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              style={{
                position: 'absolute',
                right: '15px',
                top: '50%',
                transform: 'translateY(-50%)',
                background: 'none',
                border: 'none',
                color: '#aaa',
                fontSize: '1.2rem',
                cursor: 'pointer'
              }}
            >
              {showPassword ? '🙈' : '👁️'}
            </button>
          </div>

          <button
            type="submit"
            className="btn-primary"
            disabled={loading}
            style={{
              width: '100%',
              fontSize: '16px',
              padding: '12px',
              fontWeight: '600'
            }}
          >
            {loading ? 'Регистрируем...' : 'Зарегистрироваться'}
          </button>
        </form>

        <p style={{
          textAlign: 'center',
          marginTop: '20px',
          fontSize: '14px',
          color: '#aaa'
        }}>
          Уже есть аккаунт?{' '}
          <Link to="/login" style={{ color: '#bb86fc', textDecoration: 'underline' }}>
            Войти
          </Link>
        </p>
      </div>
    </div>
  );
}