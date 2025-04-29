"use client";

import React, { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { BACKEND_URL } from "@/config";
import Header from "../../components/Header_Home";
import { User } from "@/app/types/User";
import { Portfolio } from "@/app/types/Portfolio";

// カタカナの先頭文字を英字イニシャルに変換する関数
function getInitial(kanaChar: string): string {
    const firstChar = kanaChar.slice(0, 1);
    const dictionary: Record<string, string> = {
        "ア": "A", "イ": "I", "ウ": "U", "エ": "E", "オ": "O",
        "カ": "K", "キ": "K", "ク": "K", "ケ": "K", "コ": "K",
        "サ": "S", "シ": "S", "ス": "S", "セ": "S", "ソ": "S",
        "タ": "T", "チ": "T", "ツ": "T", "テ": "T", "ト": "T",
        "ナ": "N", "ニ": "N", "ヌ": "N", "ネ": "N", "ノ": "N",
        "ハ": "H", "ヒ": "H", "フ": "H", "ヘ": "H", "ホ": "H",
        "マ": "M", "ミ": "M", "ム": "M", "メ": "M", "モ": "M",
        "ヤ": "Y", "ユ": "Y", "ヨ": "Y",
        "ラ": "R", "リ": "R", "ル": "R", "レ": "R", "ロ": "R",
        "ワ": "W", "ヲ": "W", "ン": "N",
    };
    return dictionary[firstChar] || "X";
}

// ユーザーのカタカナ名からイニシャル（例: "タニ" と "ケンシン" → "T.K."）を取得する関数
function getInitialsFromKana(firstNameKana: string, lastNameKana: string): string {
    const lastInit = getInitial(lastNameKana);
    const firstInit = getInitial(firstNameKana);
    return `${lastInit}.${firstInit}.`;
}

const PortfolioPage = () => {
    const params = useParams();
    const postId = params.id;
    const [post, setPost] = useState<Portfolio | null>(null);

    // Header 用のユーザー情報取得
    const [user, setUser] = useState<User | null>(null);
    useEffect(() => {
        fetch(`${BACKEND_URL}/user/GetInfo`, { credentials: "include" })
            .then((res) => res.json())
            .then((data) => setUser(data.user))
            .catch((error) =>
                console.error("Error fetching user info for header:", error)
            );
    }, []);

    useEffect(() => {
        if (!postId) return;
        fetch(`${BACKEND_URL}/Portfolio/${postId}`, {
            credentials: "include",
        })
            .then((res) => {
                if (!res.ok) {
                    throw new Error("Failed to fetch post");
                }
                return res.json();
            })
            .then((data) => {
                console.log(data);
                setPost(data.post);
            })
            .catch((error) => {
                console.error("Error fetching post:", error);
            });
    }, [postId]);

    if (!post) {
        return <div className="p-8">読み込み中...</div>;
    }

    // ユーザーのイニシャルを取得
    const userInitials =
        post.User.FirstNameKana && post.User.LastNameKana
            ? getInitialsFromKana(post.User.FirstNameKana, post.User.LastNameKana)
            : `${post.User.LastName} ${post.User.FirstName}`;

    return (
        <div className="min-h-screen bg-white">
            <Header
                onPostClick={() => { }}
                userHasPhoto={!!(user && user.ProfileImageURL)}
                userPhotoURL={
                    user
                        ? `${BACKEND_URL}/${user.ProfileImageURL}`
                        : "/images/defaultUserIcon.png"
                }
            />
            {/* Header 分の余白 */}
            <div className="pt-16 container mx-auto px-4 py-8">
                <div className="flex flex-col md:flex-row gap-8">
                    {/* 左カラム (画像) */}
                    <div className="md:w-1/2 w-full">
                        <div className="sticky top-0">
                            <div className="flex flex-col gap-4">
                                {post.Images && post.Images.length > 0 ? (
                                    post.Images.map((img) => (
                                        <img
                                            key={img.ID}
                                            src={`${BACKEND_URL}/${img.URL}`}
                                            alt="Post Image"
                                            className="object-cover rounded-md"
                                        />
                                    ))
                                ) : (
                                    <div>No images</div>
                                )}
                            </div>
                        </div>
                    </div>
                    {/* 右カラム (作品/ユーザー情報) */}
                    <div className="md:w-1/2 w-full flex flex-col gap-6">
                        {/* ユーザー情報 */}
                        <div className="flex items-center gap-4">
                            {post.User.ProfileImageURL ? (
                                <img
                                    src={`${BACKEND_URL}/${post.User.ProfileImageURL}`}
                                    alt="User"
                                    className="w-16 h-16 object-cover rounded-full"
                                />
                            ) : (
                                <div className="w-16 h-16 bg-gray-300 rounded-full" />
                            )}
                            <div>
                                <p className="font-bold text-xl">{userInitials}</p>
                                <p className="text-gray-500 text-sm">
                                    {post.User.SchoolName} ・ {post.User.GraduationYear}卒
                                </p>
                            </div>
                        </div>
                        {/* タイトル */}
                        <div>
                            <h1 className="text-3xl font-semibold mb-2">{post.Title}</h1>
                            <div className="flex flex-wrap gap-2">
                                {post.Genres &&
                                    post.Genres.map((genre, idx) => (
                                        <span
                                            key={idx}
                                            className="inline-block bg-orange-100 text-orange-700 px-3 py-1 rounded-full text-sm"
                                        >
                                            {genre}
                                        </span>
                                    ))}
                            </div>
                        </div>
                        {/* 説明 */}
                        <div>
                            <p className="text-gray-700 leading-relaxed">{post.Description}</p>
                        </div>
                        {/* スキル */}
                        <div>
                            <h2 className="font-bold text-lg mb-2">スキル・ツール</h2>
                            <div className="flex flex-wrap gap-2">
                                {post.Skills &&
                                    post.Skills.map((skill, idx) => (
                                        <span
                                            key={idx}
                                            className="inline-block bg-gray-200 text-gray-800 px-2 py-1 rounded text-sm"
                                        >
                                            {skill}
                                        </span>
                                    ))}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default PortfolioPage;
