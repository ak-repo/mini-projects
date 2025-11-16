export default function Card({ children, className = '', hover = false }) {
  return (
    <div
      className={`bg-[#1e293b] border border-gray-800 rounded-2xl p-6 ${
        hover ? 'hover:scale-105 hover:shadow-xl hover:shadow-[#38bdf8]/10 transition-all duration-300' : ''
      } ${className}`}
    >
      {children}
    </div>
  );
}
