// components/LoginModal.tsx
import { useRouter } from "next/navigation";
import React, { FC, useState } from "react";

interface ModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSignUpClick: () => void;
}

const LoginModal: FC<ModalProps> = ({ isOpen, onClose, onSignUpClick }) => {
    const router = useRouter();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [errors, setErrors] = useState<{ email?: string; password?: string; general?: string }>({});

    if (!isOpen) return null;

    const handleOverlayClick = () => {
        onClose();
    };

    const handleModalContentClick = (e: React.MouseEvent<HTMLDivElement>) => {
        e.stopPropagation(); // モーダル内のクリックは閉じないようにする
    };

    // メールアドレスのバリデーション関数
    const validateEmail = (email: string) => {
        const re = /\S+@\S+\.\S+/;
        return re.test(email);
    };

    const handleLogin = () => {
        let validationErrors: { email?: string; password?: string } = {};

        // ...バリデーション処理はそのまま...

        // エラーメッセージをクリア
        setErrors({});

        // 「ログイン状態を保持」のチェック状態を取得
        const isRememberMeChecked = (document.getElementById('rememberMeCheckbox') as HTMLInputElement)?.checked;

        // ログイン処理
        fetch('http://localhost:8080/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include', // クッキーを含める
            body: JSON.stringify({ email, password, rememberMe: isRememberMeChecked }),
        })
            .then(response => {
                if (!response.ok) {
                    return response.json().then(data => {
                        throw new Error(data.error || 'ログインに失敗しました');
                    }).catch(() => {
                        throw new Error('ログインに失敗しました');
                    });
                } else {
                    // レスポンスボディをパースせずに次の処理へ
                    onClose();
                    router.push("/home");
                }
            })
            .catch(error => {
                setErrors({ general: error.message });
            });
    };

    const handleGoogleLoginClick = () => {
        const isRememberMeChecked = (document.getElementById('rememberMeCheckbox') as HTMLInputElement)?.checked;

        // バックエンドのGoogleログインエンドポイントにリダイレクト
        window.location.href = `http://localhost:8080/auth/google/login?rememberMe=${isRememberMeChecked}`;
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
                <h2 className="text-xl font-bold mb-4">ログイン</h2>
                {errors.general && (
                    <div className="text-red-500 mb-4">{errors.general}</div>
                )}
                <button className="w-full py-2 mb-4 border border-gray-300 rounded-md hover:bg-gray-200 font-bold"
                    onClick={handleGoogleLoginClick}>
                    <img src="https://developers.google.com/identity/images/g-logo.png" alt="Google Icon" className="inline-block mr-2" />
                    Googleでログイン
                </button>
                <div className="text-center my-2">または</div>
                <input
                    type="email"
                    placeholder="メールアドレス"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className={`w-full mb-2 p-2 border ${errors.email ? 'border-red-500' : 'border-gray-300'} rounded-md`}
                />
                {errors.email && (
                    <div className="text-red-500 mb-2 text-sm">{errors.email}</div>
                )}
                <input
                    type="password"
                    placeholder="パスワード"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className={`w-full mb-2 p-2 border ${errors.password ? 'border-red-500' : 'border-gray-300'} rounded-md`}
                />
                {errors.password && (
                    <div className="text-red-500 mb-2 text-sm">{errors.password}</div>
                )}
                <label className="flex items-center mb-4">
                    <input type="checkbox" id="rememberMeCheckbox" className="form-checkbox" />
                    <span className="ml-2">ログイン状態を保持</span>
                </label>
                <button
                    className="bg-orange-500 w-full py-2 text-white rounded-md hover:bg-orange-600"
                    onClick={handleLogin}
                >
                    ログイン
                </button>
                <div className="text-center mt-4 text-sm">
                    パスワードをお忘れですか？
                </div>
                <div className="text-center mt-4 text-sm">
                    アカウントをお持ちですか？{" "}
                    <a
                        href="#"
                        className="text-orange-500 hover:underline"
                        onClick={onSignUpClick}
                    >
                        登録
                    </a>
                </div>
            </div>
        </div>
    );
};

export default LoginModal;
