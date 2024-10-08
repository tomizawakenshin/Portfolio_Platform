import { FC, useState } from "react";

interface EmailSignUpModalProps {
    isOpen: boolean;
    onClose: () => void;
    onComplete: (email: string) => void; // 追加
}

const EmailSignUpModal: FC<EmailSignUpModalProps> = ({ isOpen, onClose, onComplete }) => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [emailError, setEmailError] = useState("");
    const [passwordError, setPasswordError] = useState("");
    const [serverError, setServerError] = useState("");

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
            const response = await fetch("http://localhost:8080/auth/signup", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ email, password }),
            });

            if (!response.ok) {
                const data = await response.json();
                setServerError(data.error || "サインアップに失敗しました。");
            } else {
                // サインアップ成功時の処理
                onClose();          // EmailSignUpModal を閉じる
                onComplete(email);  // 親コンポーネントにサインアップ完了を通知
            }
        } catch (error) {
            console.error("Error during signup:", error);
            setServerError("サーバーエラーが発生しました。");
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
                <button
                    className="w-full py-2 bg-orange-500 text-white rounded-md hover:bg-orange-600"
                    onClick={handleSubmit}
                >
                    登録する
                </button>
                <button
                    className="w-full py-2 mt-4 text-red-500"
                    onClick={onClose}
                >
                    戻る
                </button>
            </div>
        </div>
    );
};

export default EmailSignUpModal;
