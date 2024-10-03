// components/SignUpModal.tsx
import { FC } from "react";

interface SignUpModalProps {
    isOpen: boolean;
    onClose: () => void;
    onLoginClick: () => void;  // ログインボタンをクリックしたときのハンドラを追加
}

const SignUpModal: FC<SignUpModalProps> = ({ isOpen, onClose, onLoginClick }) => {
    if (!isOpen) return null;

    const handleOverlayClick = () => {
        onClose();  // モーダルの外側をクリックした時に閉じる
    };

    const handleModalContentClick = (e: React.MouseEvent<HTMLDivElement>) => {
        e.stopPropagation();  // モーダル内をクリックした場合は閉じないようにする
    };

    return (
        <div
            className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50"
            onClick={handleOverlayClick}  // 背景のクリックで閉じる
        >
            <div
                className="bg-white p-8 rounded-lg shadow-lg max-w-lg w-full relative"
                onClick={handleModalContentClick}  // モーダル内のクリックは閉じない
            >
                <button
                    className="absolute top-4 right-4 text-gray-500 hover:text-gray-700"
                    onClick={onClose}
                >
                    &#10005;
                </button>
                <h2 className="text-xl font-bold mb-4">新規登録</h2>
                <p className="text-gray-600 mb-4">
                    このフォームではお客様の名前やメールアドレスなどの個人情報を取得いたします。これらの情報は「個人情報の取り扱いについて」に記載された利用目的の範囲内で使用いたします。
                </p>
                <div className="mb-4">
                    <label className="flex items-center space-x-2">
                        <input type="checkbox" className="form-checkbox" />
                        <span>個人情報の取り扱いについて同意する</span>
                    </label>
                </div>
                <button className="bg-red-500 text-white w-full py-2 rounded-md hover:bg-red-600">
                    Googleで登録
                </button>
                <div className="text-center my-4">または</div>
                <button className="bg-orange-500 text-white w-full py-2 rounded-md hover:bg-orange-600">
                    メールアドレスで登録する
                </button>
                <div className="text-center mt-4 text-sm">
                    アカウントをお持ちですか？{" "}
                    <a
                        href="#"
                        className="text-orange-500 hover:underline"
                        onClick={onLoginClick}  // ログインボタンのクリック時に呼び出す
                    >
                        ログイン
                    </a>
                </div>
            </div>
        </div>
    );
};

export default SignUpModal;
