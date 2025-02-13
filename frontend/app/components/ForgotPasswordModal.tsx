// components/ForgotPasswordModal.tsx
import { BACKEND_URL } from '@/config';
import React, { FC, useState } from 'react';

interface ForgotPasswordModalProps {
    isOpen: boolean;
    onClose: () => void;
}

const ForgotPasswordModal: FC<ForgotPasswordModalProps> = ({ isOpen, onClose }) => {
    const [email, setEmail] = useState('');
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');

    if (!isOpen) return null;

    const handleOverlayClick = () => {
        onClose();
    };

    const handleModalContentClick = (e: React.MouseEvent<HTMLDivElement>) => {
        e.stopPropagation();
    };

    const handleSubmit = () => {
        if (!email) {
            setError('メールアドレスを入力してください。');
            return;
        }

        // API呼び出し
        fetch(`${BACKEND_URL}/auth/RequestPasswordReset`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email }),
        })
            .then(async response => {
                const data = await response.json();
                if (!response.ok) {
                    throw new Error(data.error || 'メールの送信に失敗しました。');
                }
                setMessage(data.message || 'パスワード再設定用のリンクを送信しました。');
                setError('');
                setEmail('');
            })
            .catch(error => {
                setError(error.message);
                setMessage('');
            });
    };
    return (
        <div
            className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50"
            onClick={handleOverlayClick}
        >
            <div
                className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full relative"
                onClick={handleModalContentClick}
            >
                <button
                    className="absolute top-4 right-4 text-gray-500 hover:text-gray-700"
                    onClick={onClose}
                >
                    &#10005;
                </button>
                <h2 className="text-xl font-bold mb-4">パスワードの再設定</h2>
                <p className="mb-4">アカウントの登録メールアドレスをご入力ください。パスワードリセット用のリンクをメールいたします。</p>

                {message && <div className="text-green-500 mb-4">{message}</div>}
                {error && <div className="text-red-500 mb-4">{error}</div>}

                <input
                    type="email"
                    placeholder="メールアドレス"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full mb-4 p-2 border border-gray-300 rounded-md"
                />
                <button
                    className="bg-orange-500 w-full py-2 text-white rounded-md hover:bg-orange-600"
                    onClick={handleSubmit}
                >
                    再設定用リンクを送る
                </button>
                <button
                    className="mt-4 text-orange-500 hover:underline"
                    onClick={onClose}
                >
                    ログイン画面に戻る
                </button>
            </div>
        </div>
    );
};

export default ForgotPasswordModal;
