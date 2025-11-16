export default function Button({
  children,
  variant = 'primary',
  className = '',
  ...props
}) {
  const baseStyles = 'px-6 py-2.5 rounded-xl font-medium transition-all duration-200';

  const variants = {
    primary: 'bg-gradient-to-r from-[#38bdf8] to-[#a855f7] hover:shadow-lg hover:shadow-[#38bdf8]/30 text-white',
    secondary: 'bg-[#1e293b] border border-gray-700 hover:border-[#38bdf8] text-gray-200',
    ghost: 'hover:bg-[#1e293b] text-gray-400 hover:text-gray-200',
  };

  return (
    <button
      className={`${baseStyles} ${variants[variant]} ${className}`}
      {...props}
    >
      {children}
    </button>
  );
}
