"use client";

import React, { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { Portfolio } from "../../types/Portfolio";  // 例：エイリアスや相対パスは適宜変更
import { User } from "../../types/User";        // こちらは必要なら直接import

const PortfolioPage = () => {
    const params = useParams();
    const postId = params.id;
    const [post, setPost] = useState<Portfolio | null>(null);

    useEffect(() => {
        if (!postId) return;

        fetch(`http://localhost:8080/Portfolio/${postId}`, {
            credentials: "include",
        })
            .then((res) => {
                if (!res.ok) {
                    throw new Error("Failed to fetch post");
                }
                return res.json();
            })
            .then((data) => {
                // data.post => { ID, Title, Description, Genres, Skills, Images, User, ... }
                setPost(data.post);
            })
            .catch((error) => {
                console.error("Error fetching post:", error);
            });
    }, [postId]);

    if (!post) {
        return <div className="p-8">読み込み中...</div>;
    }

    // ここで post.User が取得できる
    return (
        <div className="min-h-screen bg-white">
            <div className="container mx-auto px-4 py-8 flex flex-col md:flex-row gap-8">

                {/* =========== 左カラム (画像) =========== */}
                <div className="md:w-1/2 w-full">
                    <div className="sticky top-0">
                        <div className="flex flex-col gap-4">
                            {post.Images.length > 0 ? (
                                post.Images.map((img) => (
                                    <img
                                        key={img.ID}
                                        src={`http://localhost:8080/${img.URL}`}
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

                {/* =========== 右カラム (作品/ユーザー情報) =========== */}
                <div className="md:w-1/2 w-full flex flex-col gap-6">
                    {/* ユーザー情報 */}
                    <div className="flex items-center gap-4">
                        {post.User.profilePictureURL ? (
                            <img
                                src={`http://localhost:8080/${post.User.profilePictureURL}`}
                                alt="User"
                                className="w-16 h-16 object-cover rounded-full"
                            />
                        ) : (
                            <div className="w-16 h-16 bg-gray-300 rounded-full" />
                        )}
                        <div>
                            <p className="font-bold text-xl">
                                {post.User.FirstName} {post.User.LastName}
                            </p>
                            <p className="text-gray-500 text-sm">{post.User.Email}</p>
                        </div>
                    </div>

                    {/* タイトル */}
                    <div>
                        <h1 className="text-3xl font-semibold mb-2">{post.Title}</h1>
                        {/* ジャンルタグ */}
                        <div className="flex flex-wrap gap-2">
                            {post.Genres.map((genre, idx) => (
                                <span
                                    key={idx}
                                    className="inline-block bg-orange-100 text-orange-700 px-3 py-1 rounded-full text-sm"
                                >
                                    {genre}
                                </span>
                            ))}
                        </div>
                    </div>

                    {/* 作品の説明 */}
                    <div>
                        <p className="text-gray-700 leading-relaxed">{post.Description}</p>
                    </div>

                    {/* スキル */}
                    <div>
                        <h2 className="font-bold text-lg mb-2">スキル・ツール</h2>
                        <div className="flex flex-wrap gap-2">
                            {post.Skills.map((skill, idx) => (
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
    );
};

export default PortfolioPage;
