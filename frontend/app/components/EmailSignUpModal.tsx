import { FC, useState } from "react";
import { useRouter } from "next/navigation";
import { BACKEND_URL } from "@/config";

interface EmailSignUpModalProps {
    isOpen: boolean;
    onClose: () => void;
    onComplete: (email: string) => void;
}

const EmailSignUpModal: FC<EmailSignUpModalProps> = ({
    isOpen,
    onClose,
    onComplete,
}) => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [emailError, setEmailError] = useState("");
    const [passwordError, setPasswordError] = useState("");
    const [serverError, setServerError] = useState("");

    // ▼追加: ローディング状態を管理
    const [isLoading, setIsLoading] = useState(false);

    const router = useRouter();

    if (!isOpen) return null;

    const handleOverlayClick = () => {
        onClose();
    };

    const handleModalContentClick = (e: React.MouseEvent<HTMLDivElement>) => {
        e.stopPropagation();
    };

    const validateEmail = (email: string) => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    };

    const handleSubmit = async () => {
        let valid = true;
        setEmailError("");
        setPasswordError("");
        setServerError("");

        if (!validateEmail(email)) {
            setEmailError("有効なメールアドレスを入力してください。");
            valid = false;
        }

        if (password.length < 8) {
            setPasswordError("パスワードは8文字以上である必要があります。");
            valid = false;
        } else if (!/^[a-zA-Z0-9]+$/.test(password)) {
            setPasswordError("パスワードは半角英数字のみを使用してください。");
            valid = false;
        }

        if (!valid) return;

        try {
            // ▼通信開始 → ローディング状態true
            setIsLoading(true);

            const response = await fetch(`${BACKEND_URL}/auth/signup`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include", // クッキーを含める
                body: JSON.stringify({ email, password }),
            });

            if (response.status === 201) {
                // アカウントが新規作成された場合の処理
                onClose();
                onComplete(email);
            } else if (response.status === 200) {
                // 既存ユーザーでログインした場合の処理
                onClose();
                router.push("/home");
            } else {
                // エラー処理
                const data = await response.json();
                setServerError(data.error || "サインアップに失敗しました。");
            }
        } catch (error) {
            console.error("Error during signup:", error);
            setServerError("サーバーエラーが発生しました。");
        } finally {
            // ▼通信終了 → ローディング状態false
            setIsLoading(false);
        }
    };

    return (
        <div
            className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50"
            onClick={handleOverlayClick}
        >
            <div
                className="bg-white p-8 rounded-lg shadow-lg max-w-lg w-full relative"
                onClick={handleModalContentClick}
            >
                <button
                    className="absolute top-4 right-4 text-gray-500 hover:text-gray-700"
                    onClick={onClose}
                >
                    &#10005;
                </button>
                <h2 className="text-xl font-bold mb-4">メールアドレスで登録</h2>
                {serverError && <p className="text-red-500 mb-4">{serverError}</p>}
                <input
                    type="email"
                    placeholder="メールアドレス"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className={`w-full p-2 mb-1 border ${emailError ? "border-red-500" : "border-gray-300"
                        } rounded`}
                />
                {emailError && (
                    <p className="text-red-500 mb-2 text-sm">{emailError}</p>
                )}
                <input
                    type="password"
                    placeholder="パスワード (8文字以上・半角英数字のみ)"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className={`w-full p-2 mb-1 border ${passwordError ? "border-red-500" : "border-gray-300"
                        } rounded`}
                />
                {passwordError && (
                    <p className="text-red-500 mb-4 text-sm">{passwordError}</p>
                )}

                {/* ▼ローディング中はボタンを押せない & テキスト変更 */}
                <button
                    className={`w-full py-2 text-white rounded-md ${isLoading
                            ? "bg-gray-400 cursor-not-allowed"
                            : "bg-orange-500 hover:bg-orange-600"
                        }`}
                    onClick={handleSubmit}
                    disabled={isLoading}
                >
                    {isLoading ? "登録中..." : "登録する"}
                </button>

                <button
                    className="w-full py-2 mt-4 text-red-500"
                    onClick={onClose}
                    disabled={isLoading}
                >
                    戻る
                </button>
            </div>
        </div>
    );
};

export default EmailSignUpModal;
