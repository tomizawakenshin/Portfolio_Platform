// components/Header.tsx
"use client";

import React from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

interface HeaderProps {
    onPostClick: () => void;
    userHasPhoto: boolean;         // ユーザーに写真があるかどうか
    userPhotoURL?: string;         // ユーザーの写真URL（あれば）
}

const Header: React.FC<HeaderProps> = ({ onPostClick, userHasPhoto, userPhotoURL }) => {
    const router = useRouter();

    const handleHomeClick = () => {
        // ホームへ移動。ルートや別のページがあれば変更
        router.push("/home");
    };

    const handleUserIconClick = () => {
        // アカウントページへ飛ぶなど
        router.push("/account");
    };

    return (
        <header className="flex items-center justify-between px-6 py-3 bg-white shadow-md">
            {/* 左側: ホーム */}
            <div>
                <button
                    onClick={handleHomeClick}
                    className="text-gray-700 hover:text-orange-500 font-semibold"
                >
                    ホーム
                </button>
            </div>

            {/* 右側: +作品を投稿 / ユーザーアイコン */}
            <div className="flex items-center space-x-4">
                <button
                    onClick={onPostClick}
                    className="
                                flex items-center 
                                space-x-2
                                px-3 py-2 
                                rounded 
                                font-bold 
                                text-white 
                                bg-orange-500 
                                hover:bg-orange-600 
                                transition-colors 
                                duration-200
                            "
                >
                    <span>+ 作品を追加</span>
                </button>

                {/* ユーザーアイコン */}
                <div
                    onClick={handleUserIconClick}
                    className="cursor-pointer relative group"
                >
                    <img
                        src={
                            userHasPhoto && userPhotoURL
                                ? userPhotoURL
                                : "/images/defaultUserIcon.png" // ← デフォルトアイコン画像へのパス
                        }
                        alt="User Icon"
                        className="w-9 h-9 rounded-full object-cover 
                       hover:opacity-90 transition-opacity duration-200"
                    />
                </div>
            </div>
        </header>
    );
};

export default Header;
