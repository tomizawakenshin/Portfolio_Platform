// components/Header.tsx
const Header = () => {
    return (
        <header className="flex justify-between items-center p-6 bg-white shadow-md">
            <div className="text-3xl font-bold text-orange-500">エンジニアのポートフォリオ</div>
            <nav className="flex space-x-8">
                <a href="#" className="text-gray-600 hover:text-orange-500">採用担当者の方へ</a>
                <a href="#" className="text-gray-600 hover:text-orange-500">学校関係者の方へ</a>
            </nav>
        </header>
    );
};

export default Header;
