// src/components/Navbar.jsx
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import Avatar from './Avatar';

export default function Navbar() {
  const token = localStorage.getItem('token');
  const email = localStorage.getItem('user_email');
  const userId = localStorage.getItem('user_id');
  const [dropdown, setDropdown] = useState(false);
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user_id');
    localStorage.removeItem('user_email');
    setDropdown(false);
    navigate('/login');
  };

  return (
    <nav style={{
      background: 'linear-gradient(90deg, #16161d, #1a1a24)',
      padding: '15px 30px',
      display: 'flex',
      justifyContent: 'space-between',
      alignItems: 'center',
      borderBottom: '1px solid #2d2d3a',
      boxShadow: '0 2px 10px rgba(0, 0, 0, 0.3)'
    }}>
      <div>
        <Link to="/" style={{
          color: '#bb86fc',
          fontSize: '1.5rem',
          fontWeight: 'bold',
          textDecoration: 'none',
          background: 'linear-gradient(90deg, #bb86fc, #03dac6)',
          WebkitBackgroundClip: 'text',
          WebkitTextFillColor: 'transparent'
        }}>
          CalcPro
        </Link>
      </div>

      <div style={{ display: 'flex', gap: '20px', alignItems: 'center' }}>
        {token ? (
          <>
            <NavLink to="/calc">Калькулятор</NavLink>
            <NavLink to="/examples">История</NavLink>
            <NavLink to="/about">О проекте</NavLink>

            {/* Иконка пользователя */}
            <div style={{ position: 'relative' }}>
              <Avatar 
                email={email} 
                onClick={() => setDropdown(!dropdown)} 
              />

              {/* Выпадающее меню */}
              {dropdown && (
                <div style={{
                  position: 'absolute',
                  top: '50px',
                  right: 0,
                  backgroundColor: '#1a1a24',
                  borderRadius: '12px',
                  boxShadow: '0 10px 30px rgba(0,0,0,0.4)',
                  width: '400px',
                  zIndex: 1000,
                  overflow: 'hidden',
                  border: '1px solid #2d2d3a'
                }}>
                  <div style={{ 
                    padding: '20px', 
                    borderBottom: '1px solid #2d2d3a',
                    background: 'linear-gradient(135deg, #1e1e24, #16161d)'
                  }}>
                    <h4 style={{ 
                      margin: '0 0 10px 0', 
                      color: '#bb86fc',
                      fontSize: '1.1rem'
                    }}>
                      Профиль
                    </h4>
                    <div style={{ 
                      display: 'flex', 
                      alignItems: 'flex-start', 
                      gap: '12px',
                      marginBottom: '15px'
                    }}>
                      <div style={{
                        width: '40px',
                        height: '40px',
                        borderRadius: '50%',
                        background: 'linear-gradient(135deg, #bb86fc, #673ab7)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: '18px',
                        fontWeight: 'bold',
                        color: 'white'
                      }}>
                        {email ? email.charAt(0).toUpperCase() : '?'}
                      </div>
                      <div>
                        <div style={{ 
                          color: '#e0e0e0', 
                          fontSize: '0.95rem',
                          wordBreak: 'break-all',
                          maxWidth: '300px'
                        }}>
                          <strong>{email || '—'}</strong>
                        </div>
                        <div style={{ 
                          color: '#aaa', 
                          fontSize: '0.8rem',
                          marginTop: '4px'
                        }}>
                          Пользователь
                        </div>
                      </div>
                    </div>
                  </div>

                  <div style={{ padding: '15px 20px' }}>
                    <div style={{ 
                      display: 'flex', 
                      alignItems: 'flex-start', 
                      gap: '10px',
                      marginBottom: '10px'
                    }}>
                      <span style={{ color: '#03dac6' }}>🆔</span>
                      <div>
                        <div style={{ 
                          color: '#aaa', 
                          fontSize: '0.8rem',
                          marginBottom: '4px'
                        }}>
                          User ID
                        </div>
                        <code style={{
                          background: '#2d2d3a',
                          padding: '8px 12px',
                          borderRadius: '6px',
                          fontSize: '0.85rem',
                          color: '#e0e0e0',
                          fontFamily: 'monospace',
                          wordBreak: 'break-all',
                          lineHeight: '1.4',
                          maxWidth: '320px'
                        }}>
                          {userId || '—'}
                        </code>
                      </div>
                    </div>
                  </div>

                  <button
                    onClick={handleLogout}
                    style={{
                      width: '100%',
                      padding: '15px 20px',
                      background: 'linear-gradient(90deg, #cf6679, #b00020)',
                      color: 'white',
                      border: 'none',
                      textAlign: 'center',
                      cursor: 'pointer',
                      fontSize: '1rem',
                      fontWeight: '500',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      gap: '8px'
                    }}
                  >
                    🔌 Выйти
                  </button>
                </div>
              )}
            </div>
          </>
        ) : (
          <>
            <NavLink to="/about">О проекте</NavLink>
            <NavLink to="/login">Вход</NavLink>
            <Link to="/register" style={{
              ...linkStyle,
              background: 'linear-gradient(135deg, #bb86fc, #673ab7)',
              color: 'white',
              padding: '8px 20px',
              borderRadius: '6px'
            }}>Регистрация</Link>
          </>
        )}
      </div>

      {/* Закрытие меню при клике вне его */}
      {dropdown && (
        <div
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            zIndex: 999,
            cursor: 'default'
          }}
          onClick={() => setDropdown(false)}
        />
      )}
    </nav>
  );
}

// ✅ Новый компонент для стилизованных ссылок
function NavLink({ to, children }) {
  return (
    <Link 
      to={to} 
      style={{
        color: '#bb86fc',
        textDecoration: 'none',
        padding: '8px 15px',
        borderRadius: '6px',
        transition: 'all 0.3s ease',
        fontWeight: '500',
        position: 'relative',
        background: 'transparent'
      }}
      onMouseEnter={(e) => {
        e.target.style.background = 'rgba(187, 134, 252, 0.1)';
        e.target.style.transform = 'translateY(-1px)';
      }}
      onMouseLeave={(e) => {
        e.target.style.background = 'transparent';
        e.target.style.transform = 'translateY(0)';
      }}
    >
      {children}
    </Link>
  );
}

// Стиль для обычных ссылок (регистрация)
const linkStyle = {
  color: '#e0e0e0',
  textDecoration: 'none',
  padding: '8px 15px',
  borderRadius: '6px',
  transition: 'all 0.3s ease',
  fontWeight: '500'
};