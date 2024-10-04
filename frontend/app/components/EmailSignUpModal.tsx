import { FC } from "react";

interface EmailSignUpModalProps {
    isOpen: boolean;
    onClose: () => void;
}

const EmailSignUpModal: FC<EmailSignUpModalProps> = ({ isOpen, onClose }) => {
    if (!isOpen) return null;

    const handleOverlayClick = () => {
        onClose();
    };

    const handleModalContentClick = (e: React.MouseEvent<HTMLDivElement>) => {
        e.stopPropagation();
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
                <input
                    type="email"
                    placeholder="メールアドレス"
                    className="w-full p-2 mb-4 border border-gray-300 rounded"
                />
                <input
                    type="password"
                    placeholder="パスワード (8文字以上・半角英数字のみ)"
                    className="w-full p-2 mb-4 border border-gray-300 rounded"
                />
                <button className="w-full py-2 bg-orange-500 text-white rounded-md hover:bg-orange-600">
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
