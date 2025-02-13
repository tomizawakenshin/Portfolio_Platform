"use client";

import React, { useEffect, useRef, useState } from "react";
import { User } from "../types/User";
import { BACKEND_URL } from "@/config";

interface ProfileEditModalProps {
    isOpen: boolean;
    onClose: () => void;
    user: User;                         // 現在のユーザー情報
    onSaveSuccess: (updatedUser: User) => void;
    // 保存が成功したら、更新された User を親に返す
}

const ProfileEditModal: React.FC<ProfileEditModalProps> = ({
    isOpen,
    onClose,
    user,
    onSaveSuccess
}) => {
    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");
    const [firstNameKana, setFirstNameKana] = useState("");
    const [lastNameKana, setLastNameKana] = useState("");

    // ユーザーアイコン用ファイル
    const [iconFile, setIconFile] = useState<File | null>(null);
    const [iconPreview, setIconPreview] = useState<string>("");

    // モーダルが開くたびに初期化
    useEffect(() => {
        if (isOpen && user) {
            setFirstName(user.FirstName);
            setLastName(user.LastName);
            setFirstNameKana(user.FirstNameKana);
            setLastNameKana(user.LastNameKana);

            // プレビューは、現在のProfileImageURLを表示するか、
            // あるいは未変更なら何もしない
            if (user.ProfileImageURL) {
                setIconPreview(`${BACKEND_URL}/${user.ProfileImageURL}`);
            } else {
                setIconPreview("");
            }
            setIconFile(null);
        }
    }, [isOpen, user]);

    const handleIconChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            setIconFile(file);
            // プレビュー用URL
            const previewUrl = URL.createObjectURL(file);
            setIconPreview(previewUrl);
        }
    };

    const handleSave = async () => {
        try {
            // multipart/form-data で送信
            const formData = new FormData();
            // テキスト部分
            formData.append("firstName", firstName);
            formData.append("lastName", lastName);
            formData.append("firstNameKana", firstNameKana);
            formData.append("lastNameKana", lastNameKana);

            // 画像が選ばれていれば追加
            if (iconFile) {
                formData.append("profileImage", iconFile);
                // コントローラ側で ctx.MultipartForm().File["profileImage"] で受け取る想定
            }

            // ここでは同じ UpdateMinimumUserInfo に PUT するが
            // 例えばPOSTにしても構わない。Golang側で ParseMultipartForm() を呼ぶ必要がある
            const res = await fetch(`${BACKEND_URL}/user/UpdateMinimumUserInfo`, {
                method: "PUT",
                credentials: "include",
                body: formData,
            });
            if (!res.ok) {
                throw new Error("プロフィール更新に失敗しました");
            }
            const data = await res.json();
            onSaveSuccess(data.user);
            onClose();
        } catch (error) {
            console.error(error);
            alert("プロフィールの更新中にエラーが発生しました。");
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white w-full max-w-2xl rounded p-6 relative">
                <button
                    className="absolute top-2 right-2 text-gray-500 hover:text-gray-800"
                    onClick={onClose}
                >
                    &times;
                </button>
                <h2 className="text-xl font-bold mb-4">プロフィールの編集</h2>

                {/* アイコンのプレビュー */}
                <div className="flex items-center justify-center mb-4">
                    <label className="relative">
                        <img
                            src={iconPreview || "/images/defaultUserIcon.png"}
                            alt="icon preview"
                            className="w-24 h-24 object-cover rounded-full border-2 border-gray-300"
                        />
                        <input
                            type="file"
                            accept="image/*"
                            className="hidden"
                            onChange={handleIconChange}
                        />
                        <span className="absolute inset-0 flex items-center justify-center text-white bg-black bg-opacity-30 opacity-0 hover:opacity-100 transition-opacity rounded-full cursor-pointer">
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                className="h-6 w-6"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                                strokeWidth={2}
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    d="M5.121 17.804A3 3 0 016 17h12a3 3 0 012.995 2.824L21 20v1a1 1 0 01-1 1H4a1 1 0 01-1-1v-1l.005-.176A3 3 0 016 17h0m-4-5a9 9 0 1118 0H2z"
                                />
                            </svg>
                        </span>
                    </label>
                </div>

                {/* 名前欄 */}
                <div className="mb-4 flex space-x-4">
                    <div className="w-1/2">
                        <label className="block text-sm font-semibold mb-1">姓</label>
                        <input
                            type="text"
                            value={lastName}
                            onChange={(e) => setLastName(e.target.value)}
                            className="border p-2 w-full rounded"
                        />
                    </div>
                    <div className="w-1/2">
                        <label className="block text-sm font-semibold mb-1">名</label>
                        <input
                            type="text"
                            value={firstName}
                            onChange={(e) => setFirstName(e.target.value)}
                            className="border p-2 w-full rounded"
                        />
                    </div>
                </div>

                {/* フリガナ */}
                <div className="mb-4 flex space-x-4">
                    <div className="w-1/2">
                        <label className="block text-sm font-semibold mb-1">セイ</label>
                        <input
                            type="text"
                            value={lastNameKana}
                            onChange={(e) => setLastNameKana(e.target.value)}
                            className="border p-2 w-full rounded"
                        />
                    </div>
                    <div className="w-1/2">
                        <label className="block text-sm font-semibold mb-1">メイ</label>
                        <input
                            type="text"
                            value={firstNameKana}
                            onChange={(e) => setFirstNameKana(e.target.value)}
                            className="border p-2 w-full rounded"
                        />
                    </div>
                </div>

                {/* 保存ボタン */}
                <div className="flex justify-end mt-4">
                    <button
                        className="px-4 py-2 bg-orange-500 text-white rounded hover:bg-orange-600"
                        onClick={handleSave}
                    >
                        保存する
                    </button>
                </div>
            </div>
        </div>
    );
};

export default ProfileEditModal;
