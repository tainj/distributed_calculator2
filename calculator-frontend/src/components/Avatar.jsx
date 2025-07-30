// src/components/Avatar.jsx
export default function Avatar({ email, size = 40, onClick }) {
  if (!email) return null;

  const letter = email[0].toUpperCase();
  const color = stringToColor(email); // Генерируем цвет по email

  return (
    <button
      onClick={onClick}
      style={{
        width: size + 'px',
        height: size + 'px',
        borderRadius: '50%',
        backgroundColor: color,
        color: 'white',
        fontWeight: 'bold',
        fontSize: (size * 0.45) + 'px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        border: 'none',
        cursor: 'pointer',
        boxShadow: '0 2px 8px rgba(0,0,0,0.2)',
        marginLeft: '10px'
      }}
    >
      {letter}
    </button>
  );
}

// Функция: генерирует цвет на основе email (чтобы был постоянный)
function stringToColor(str) {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }
  const hue = hash % 360;
  return `hsl(${hue}, 70%, 50%)`; // Фиолетовые/синие оттенки
}