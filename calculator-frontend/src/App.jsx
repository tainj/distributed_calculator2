// src/App.js
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Navbar from './components/Navbar';
import Home from './components/Home';
import Register from './components/Register';
import Login from './components/Login';
import Calculator from './components/Calculator';
import Examples from './components/Examples';
import About from './components/About';

function App() {
  return (
    <Router>
      <Navbar />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/register" element={<Register />} />
        <Route path="/login" element={<Login />} />
        <Route path="/calc" element={<Calculator />} />
        <Route path="/examples" element={<Examples />} />
        <Route path="/about" element={<About />} />
      </Routes>
    </Router>
  );
}

export default App;