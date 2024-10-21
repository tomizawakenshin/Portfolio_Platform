import { useRouter } from 'next/navigation';
import React, { useState } from 'react';

interface MnimumUserInfoInputModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSubmit: (firstName: string, lastName: string) => void;
}

const MnimumUserInfoInputModal: React.FC<MnimumUserInfoInputModalProps> = ({ isOpen, onClose, onSubmit }) => {
    const router = useRouter();

    const [firstName, setFirstName] = useState('');
    const [lastName, setLastName] = useState('');

    if (!isOpen) return null;

    const handleSubmit = () => {
        onSubmit(firstName, lastName);
        onClose();
    };

    const handleLogout = () => {
        fetch('http://localhost:8080/auth/logout', {
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
        <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
            {/* モーダル全体を包む relative なコンテナを追加 */}
            <div className="relative">
                <div className="bg-white p-8 rounded-lg shadow-lg max-w-lg w-full">
                    <h2 className="text-xl font-bold mb-4">基本情報の登録</h2>
                    <p className="text-gray-700 mb-4">
                        サービスをご利用いただくにあたって、以下の情報を入力してください。
                    </p>
                    {/* 名前 必須のラベル */}
                    <div className="mb-4">
                        <label className="text-gray-700 font-bold inline-block mb-2">
                            名前 <span className="text-orange-500">必須</span>
                        </label>
                        {/* 姓と名を横並びに配置 */}
                        <div className="flex space-x-4">
                            <input
                                type="text"
                                placeholder="姓"
                                value={lastName}
                                onChange={(e) => setLastName(e.target.value)}
                                className="w-1/2 p-2 border border-gray-300 rounded"
                            />
                            <input
                                type="text"
                                placeholder="名"
                                value={firstName}
                                onChange={(e) => setFirstName(e.target.value)}
                                className="w-1/2 p-2 border border-gray-300 rounded"
                            />
                        </div>
                    </div>
                    <button
                        className="w-full py-2 bg-orange-500 text-white rounded-md hover:bg-orange-600"
                        onClick={handleSubmit}
                    >
                        登録する
                    </button>
                </div>
                {/* ログアウトボタンをモーダルボディの外側に配置 */}
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

export default MnimumUserInfoInputModal;
