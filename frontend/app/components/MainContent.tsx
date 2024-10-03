// components/MainContent.tsx
const MainContent = () => {
    return (
        <section className="flex flex-col items-center justify-center min-h-screen bg-gray-50">
            <h1 className="text-5xl font-bold text-gray-900">タイトル（仮）</h1>
            <p className="mt-4 text-xl text-gray-600 text-center">
                このサイトに関しての簡単な説明（仮）
            </p>
            <div className="mt-8 flex space-x-4">
                <button className="px-6 py-3 bg-orange-500 text-white rounded-lg text-lg font-bold hover:bg-orange-600">
                    無料ではじめる
                </button>
                <button className="px-6 py-3 border border-orange-500 text-orange-500 rounded-lg text-lg font-bold hover:bg-orange-100">
                    ログイン
                </button>
            </div>
        </section>
    );
};

export default MainContent;
