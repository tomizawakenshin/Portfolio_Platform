import { BACKEND_URL } from "@/config";
import { useRouter } from "next/navigation";
import React, { FC, useState } from "react";

interface ModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSignUpClick: () => void;
    onForgotPasswordClick: () => void;
}

const LoginModal: FC<ModalProps> = ({
    isOpen,
    onClose,
    onSignUpClick,
    onForgotPasswordClick,
}) => {
    const router = useRouter();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [errors, setErrors] = useState<{
        email?: string;
        password?: string;
        general?: string;
    }>({});
    const [isForgotPasswordModalOpen, setIsForgotPasswordModalOpen] =
        useState(false);

    // ▼追加: ローディング状態を管理
    const [isLoading, setIsLoading] = useState(false);

    if (!isOpen) return null;

    const handleOverlayClick = () => {
        onClose();
    };

    const handleModalContentClick = (e: React.MouseEvent<HTMLDivElement>) => {
        e.stopPropagation();
    };

    const validateEmail = (email: string) => {
        const re = /\S+@\S+\.\S+/;
        return re.test(email);
    };

    const handleLogin = async () => {
        let validationErrors: { email?: string; password?: string } = {};
        // エラーメッセージをクリア
        setErrors({});

        const isRememberMeChecked = (
            document.getElementById("rememberMeCheckbox") as HTMLInputElement
        )?.checked;

        // 簡易バリデーション例
        if (!validateEmail(email)) {
            validationErrors.email = "有効なメールアドレスを入力してください。";
        }
        if (!password) {
            validationErrors.password = "パスワードを入力してください。";
        }

        if (Object.keys(validationErrors).length > 0) {
            setErrors(validationErrors);
            return;
        }

        try {
            // ▼ローディング開始
            setIsLoading(true);

            const response = await fetch(`${BACKEND_URL}/auth/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include", // クッキーを含める
                body: JSON.stringify({ email, password, rememberMe: isRememberMeChecked }),
            });

            if (!response.ok) {
                const data = await response.json().catch(() => { });
                throw new Error(data?.error || "ログインに失敗しました");
            }

            // 成功時
            onClose();
            router.push("/home");
        } catch (error: any) {
            setErrors({ general: error.message });
            console.error("Login error:", error);
        } finally {
            // ▼ローディング終了
            setIsLoading(false);
        }
    };

    const handleGoogleLoginClick = () => {
        const isRememberMeChecked = (
            document.getElementById("rememberMeCheckbox") as HTMLInputElement
        )?.checked;

        // バックエンドのGoogleログインエンドポイントにリダイレクト
        window.location.href = `${BACKEND_URL}/auth/google/login?rememberMe=${isRememberMeChecked}`;
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
                    disabled={isLoading}
                >
                    &#10005;
                </button>
                <h2 className="text-xl font-bold mb-4">ログイン</h2>
                {errors.general && (
                    <div className="text-red-500 mb-4">{errors.general}</div>
                )}

                <button
                    className={`w-full py-2 mb-4 border border-gray-300 rounded-md font-bold ${isLoading
                            ? "cursor-not-allowed bg-gray-100 text-gray-500"
                            : "hover:bg-gray-200"
                        }`}
                    onClick={handleGoogleLoginClick}
                    disabled={isLoading}
                >
                    <img
                        src="https://developers.google.com/identity/images/g-logo.png"
                        alt="Google Icon"
                        className="inline-block mr-2"
                    />
                    Googleでログイン
                </button>

                <div className="text-center my-2">または</div>

                <input
                    type="email"
                    placeholder="メールアドレス"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className={`w-full mb-2 p-2 border ${errors.email ? "border-red-500" : "border-gray-300"
                        } rounded-md`}
                    disabled={isLoading}
                />
                {errors.email && (
                    <div className="text-red-500 mb-2 text-sm">{errors.email}</div>
                )}

                <input
                    type="password"
                    placeholder="パスワード"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className={`w-full mb-2 p-2 border ${errors.password ? "border-red-500" : "border-gray-300"
                        } rounded-md`}
                    disabled={isLoading}
                />
                {errors.password && (
                    <div className="text-red-500 mb-2 text-sm">{errors.password}</div>
                )}

                <label className="flex items-center mb-4">
                    <input
                        type="checkbox"
                        id="rememberMeCheckbox"
                        className="form-checkbox"
                        disabled={isLoading}
                    />
                    <span className="ml-2">ログイン状態を保持</span>
                </label>

                <button
                    className={`w-full py-2 text-white rounded-md ${isLoading
                            ? "bg-gray-400 cursor-not-allowed"
                            : "bg-orange-500 hover:bg-orange-600"
                        }`}
                    onClick={handleLogin}
                    disabled={isLoading}
                >
                    {isLoading ? "ログイン中..." : "ログイン"}
                </button>

                <div
                    className="text-center mt-4 text-sm text-orange-500 hover:underline cursor-pointer"
                    onClick={() => {
                        onForgotPasswordClick();
                    }}
                >
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
