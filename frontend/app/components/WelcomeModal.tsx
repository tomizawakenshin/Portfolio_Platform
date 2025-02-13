import { BACKEND_URL } from '@/config';
import { useRouter } from 'next/navigation';
import React from 'react';

interface WelcomeModalProps {
    isOpen: boolean;
    onClose: () => void;
    onStart: () => void;
}

const WelcomeModal: React.FC<WelcomeModalProps> = ({ isOpen, onClose, onStart }) => {
    const router = useRouter();

    if (!isOpen) {
        return null;
    }

    const handleLogout = () => {
        fetch(`${BACKEND_URL}/auth/logout`, {
            method: 'POST',
            credentials: 'include',
        })
            .then(response => {
                if (response.ok) {
                    // ログアウト成功時の処理
                    router.push("/auth") // ログインページにリダイレクト
                } else {
                    throw new Error('Failed to logout');
                }
            })
            .catch(error => {
                console.error('Error during logout:', error);
            });
    };
    return (
        <div className="fixed inset-0 bg-orange-500 bg-opacity-50 flex justify-center items-center z-50">
            <div className='relative'>
                <div className="bg-white rounded-lg shadow-lg max-w-xl w-full px-6 py-24 text-center">
                    <div className="text-6xl mb-4">⭐️⭐️⭐️</div>
                    <h2 className="text-2xl font-bold mb-4 text-orange-500">
                        ようこそ、タイトル名（仮）へ
                    </h2>
                    <p className="text-gray-600 mb-6">
                        学生エンジニアのためのポートフォリオプラットフォームです。
                        自分のポートフォリオを公開してフィードバックをもらったり、他の学生エンジニアのポートフォリオを参考にしたりできます。
                        同世代の仲間とつながるきっかけを作りましょう。
                    </p>
                    <button
                        className="bg-orange-500 text-white px-6 py-2 rounded-lg hover:bg-orange-600 transition-colors"
                        onClick={onStart} // 修正
                    >
                        はじめる
                    </button>
                </div>
                <button
                    className="absolute -bottom-10 right-0 text-white px-4 py-2 rounded-md hover:text-gray-300"
                    onClick={handleLogout}
                >
                    ログアウト
                </button>
            </div>
        </div>
    );
};

export default WelcomeModal;
