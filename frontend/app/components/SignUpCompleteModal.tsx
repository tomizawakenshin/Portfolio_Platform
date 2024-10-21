import { FC } from "react";

interface SignUpCompleteModalProps {
    isOpen: boolean;
    onClose: () => void;
    email: string;
}

const SignUpCompleteModal: FC<SignUpCompleteModalProps> = ({ isOpen, onClose, email }) => {
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
                <div className="text-center">
                    <div className="mb-4">
                        <span className="text-orange-500 text-4xl">✔</span>
                    </div>
                    <h2 className="text-xl font-bold mb-4">仮登録を受け付けました</h2>
                    <p className="text-gray-700 mb-2">
                        メールアドレスに送られたURLから本登録を行ってください。
                    </p>
                    <p className="font-bold text-lg">{email}</p>
                    <p className="text-gray-600 mt-4">に送信しました。</p>
                </div>
            </div>
        </div>
    );
};

export default SignUpCompleteModal;
