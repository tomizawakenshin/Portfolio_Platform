// components/LoginModal.tsx
import { FC } from "react";

interface ModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSignUpClick: () => void;
}

const LoginModal: FC<ModalProps> = ({ isOpen, onClose, onSignUpClick }) => {
    if (!isOpen) return null;

    const handleOverlayClick = () => {
        onClose();
    };

    const handleModalContentClick = (e: React.MouseEvent<HTMLDivElement>) => {
        e.stopPropagation();  // モーダル内のクリックは閉じないようにする
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
                <button className=" w-full py-2 mb-4 border border-gray-300 rounded-md hover:bg-gray-200 font-bold">
                    <img src="https://developers.google.com/identity/images/g-logo.png" alt="Google Icon" className="inline-block mr-2" />
                    Googleでログイン
                </button>
                <div className="text-center my-2">または</div>
                <input
                    type="email"
                    placeholder="メールアドレス"
                    className="w-full mb-4 p-2 border border-gray-300 rounded-md"
                />
                <input
                    type="password"
                    placeholder="パスワード"
                    className="w-full mb-4 p-2 border border-gray-300 rounded-md"
                />
                <label className="flex items-center mb-4">
                    <input type="checkbox" className="form-checkbox" />
                    <span className="ml-2">ログイン状態を保持</span>
                </label>
                <button className="bg-orange-500 w-full py-2 text-white rounded-md hover:bg-orange-600">
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
                        onClick={onSignUpClick}  // ログインボタンのクリック時に呼び出す
                    >
                        登録
                    </a>
                </div>
            </div>
        </div>
    );
};

export default LoginModal;
