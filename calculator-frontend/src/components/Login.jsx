// src/components/Login.jsx
import { useState } from 'react';
import api from '../services/api';
import { useNavigate } from 'react-router-dom';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const res = await api.post('/v1/login', { email, password });

      console.log('Login response:', res.data);

      if (res.data.success) {
        localStorage.setItem('token', res.data.token);
        localStorage.setItem('user_id', res.data.userId);
        localStorage.setItem('user_email', email);
        navigate('/');
      } else {
        setError(res.data.error || 'Неверный email или пароль');
      }
    } catch (err) {
      setError('Не удалось подключиться к серверу');
      console.error('Login error:', err);
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
          Вход в CalcPro
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

          <div style={{ position: 'relative', marginBottom: '20px' }}>
            <span style={{
              position: 'absolute',
              left: '15px',
              top: '50%',
              transform: 'translateY(-50%)',
              color: '#aaa'
            }}>
              🔒
            </span>
            <input
              className="input-field"
              type="password"
              placeholder="Пароль"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              style={{ paddingLeft: '45px' }}
            />
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
            {loading ? 'Входим...' : 'Войти'}
          </button>
        </form>

        <p style={{
          textAlign: 'center',
          marginTop: '20px',
          fontSize: '14px',
          color: '#aaa'
        }}>
          Нет аккаунта?{' '}
          <Link to="/register" style={{ color: '#bb86fc', textDecoration: 'underline' }}>
            Зарегистрироваться
          </Link>
        </p>
      </div>
    </div>
  );
}

// Не забудь импортировать Link
import { Link } from 'react-router-dom';
``