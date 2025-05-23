// components/SignUpModal.tsx
import { BACKEND_URL } from "@/config";
import { useRouter } from "next/navigation";
import { FC, useState } from "react";

interface SignUpModalProps {
    isOpen: boolean;
    onClose: () => void;
    onLoginClick: () => void; // ログインボタンをクリックしたときのハンドラ
    onEmailSignUpClick: () => void; // メールアドレスで登録するボタンをクリックしたときのハンドラ
}

const SignUpModal: FC<SignUpModalProps> = ({ isOpen, onClose, onLoginClick, onEmailSignUpClick }) => {
    const router = useRouter();
    const [isAgreed, setIsAgreed] = useState(false);

    if (!isOpen) return null;

    const handleOverlayClick = () => {
        onClose();
    };

    const handleModalContentClick = (e: React.MouseEvent<HTMLDivElement>) => {
        e.stopPropagation();
    };

    const handleGoogleSignUpClick = () => {
        if (isAgreed) {
            // バックエンドのGoogleログインエンドポイントにリダイレクト
            window.location.href = `${BACKEND_URL}/auth/google/login`;
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
                <h2 className="text-xl font-bold mb-4">新規登録</h2>
                <div className="border border-gray-300 p-4 rounded mb-4">
                    <p className="text-gray-600 mb-4">
                        このフォームでは、ご利用に必要な最小限の個人情報のみをご入力いただきます。<br />
                        入力いただいた情報は、厳重に管理されます。<br />
                        ご登録の際は、上記内容にご同意いただいた上で送信してください。
                    </p>
                </div>
                <div className="flex items-center mb-4 cursor-pointer">
                    <label className="flex items-center">
                        <input
                            type="checkbox"
                            className="form-checkbox cursor-pointer"
                            checked={isAgreed}
                            onChange={(e) => setIsAgreed(e.target.checked)}
                        />
                        <span className="ml-2 text-sm cursor-pointer">個人情報の取り扱いについて同意する</span>
                    </label>
                </div>
                <button
                    className={`w-full py-2 rounded-md flex items-center justify-center mb-4 ${isAgreed ? "bg-red-500 hover:bg-red-600" : "bg-gray-400 cursor-not-allowed"} text-white`}
                    disabled={!isAgreed}
                    onClick={handleGoogleSignUpClick}
                >
                    <span className="mr-2"></span>Googleで登録
                </button>
                <div className="text-center mb-4 text-gray-600">または</div>
                <button
                    className={`text-white w-full py-2 rounded-md ${isAgreed ? "bg-orange-500 hover:bg-orange-600" : "bg-gray-400 cursor-not-allowed"}`}
                    onClick={onEmailSignUpClick} // メールアドレスで登録するボタンのクリック時に呼び出す
                >
                    メールアドレスで登録する
                </button>
                <div className="text-center mt-4 text-sm">
                    アカウントをお持ちですか？{" "}
                    <a
                        href="#"
                        className="text-orange-500 hover:underline"
                        onClick={onLoginClick} // ログインボタンのクリック時に呼び出す
                    >
                        ログイン
                    </a>
                </div>
            </div>
        </div>
    );
};

export default SignUpModal;
