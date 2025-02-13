'use client'

import React, { useState, useEffect } from "react";
import { useRouter } from 'next/navigation';
import { BACKEND_URL } from "@/config";

interface PasswordResetProps {
    params: {
        token: string;
    };
}

const PasswordReset: React.FC<PasswordResetProps> = ({ params }) => {
    const router = useRouter();
    const { token } = params;
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(true); // ローディング状態を追加

    useEffect(() => {
        const validateToken = async () => {
            try {
                const response = await fetch(`${BACKEND_URL}/auth/CheckResetToken`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ token }),
                });

                if (!response.ok) {
                    throw new Error('トークンが無効または期限切れです。');
                }

                setLoading(false); // トークンが有効であることを確認
            } catch (err) {
                // トークンが無効または期限切れの場合
                router.push('/auth'); // 指定のURLにリダイレクト
            }
        };

        if (token) {
            validateToken();
        } else {
            router.push('/auth'); // トークンがない場合もリダイレクト
        }
    }, [token, router]);

    const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setPassword(e.target.value);
    };

    const handleConfirmPasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setConfirmPassword(e.target.value);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (password !== confirmPassword) {
            alert("パスワードが一致しません。");
            return;
        }
        if (!token) {
            alert("無効なトークンです。");
            return;
        }

        try {
            const response = await fetch(`${BACKEND_URL}/auth/ResetPassword`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include', // クッキーを含める
                body: JSON.stringify({ token, new_password: password }),
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'パスワードのリセットに失敗しました。');
            }

            router.push('/home'); // ホームページにリダイレクト
        } catch (err: any) {
            // エラーメッセージを表示
            alert(err.message);
        }
    };

    if (loading) {
        return <div>読み込み中...</div>; // トークン検証中の表示
    }

    return (
        <div className="flex justify-center items-center h-screen">
            <div className="w-full max-w-md bg-white p-8 ">
                <h2 className="text-2xl font-semibold text-center mb-6">パスワードの再設定</h2>
                <form onSubmit={handleSubmit}>
                    <div className="mb-4">
                        <label className="block text-gray-700">新しいパスワード（8文字以上・半角英数字のみ）</label>
                        <input
                            type="password"
                            className="w-full mt-1 p-2 border border-gray-300 rounded-md focus:outline-none focus:border-orange-400"
                            value={password}
                            onChange={handlePasswordChange}
                            required
                            minLength={8}
                        />
                    </div>
                    <div className="mb-6">
                        <label className="block text-gray-700">新しいパスワードの確認</label>
                        <input
                            type="password"
                            className="w-full mt-1 p-2 border border-gray-300 rounded-md focus:outline-none focus:border-orange-400"
                            value={confirmPassword}
                            onChange={handleConfirmPasswordChange}
                            required
                        />
                    </div>
                    <div className="flex justify-between items-center">
                        <button
                            type="button"
                            className="bg-gray-300 text-gray-700 py-2 px-4 rounded-md hover:bg-gray-400"
                            onClick={() => router.push('/login')}
                        >
                            キャンセル
                        </button>
                        <button
                            type="submit"
                            className="bg-orange-500 text-white py-2 px-4 rounded-md hover:bg-orange-600"
                        >
                            更新
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default PasswordReset;
