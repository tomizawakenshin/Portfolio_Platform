"use client";

import { BACKEND_URL } from "@/config";
import { useSearchParams, useRouter } from "next/navigation";
import { useState, useEffect } from "react";

export default function VerifyStartPage() {
    const searchParams = useSearchParams();
    const router = useRouter();
    const token = searchParams.get("token");

    const [message, setMessage] = useState("");
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        if (!token) {
            setMessage("トークンがありません。URLを確認してください。");
        }
    }, [token]);

    const handleVerify = async () => {
        if (!token) return;
        setLoading(true);
        setMessage("");

        try {
            const res = await fetch(`${BACKEND_URL}/auth/verify`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include", // Cookieを受け取る
                body: JSON.stringify({ token }),
            });
            if (!res.ok) {
                // レスポンスが失敗ステータスなら詳細を取得
                const data = await res.json().catch(() => ({}));
                throw new Error(data.error || "本登録に失敗しました");
            }
            const data = await res.json();
            setMessage(data.message || "本登録が完了しました！");

            // 登録完了後に /home へ移動する
            router.push("/home");
        } catch (err: any) {
            setMessage(err.message);
        } finally {
            setLoading(false);
        }
    };

    if (!token) {
        return (
            <div className="flex items-center justify-center h-screen p-4">
                <div className="text-red-500 text-center">
                    {message}
                </div>
            </div>
        );
    }

    return (
        <div className="flex items-center justify-center h-screen p-4 bg-gray-50">
            <div className="bg-white shadow-md rounded-lg p-8 w-full max-w-md">
                <h1 className="text-2xl font-bold mb-4 text-center">本登録の確認</h1>

                <p className="mb-6 text-gray-700 text-center">
                    メール内リンクを踏んでいただきありがとうございます。
                    <br />
                    下記ボタンを押すと本登録が完了します。
                </p>

                {/* ボタン */}
                <button
                    className={`w-full py-2 text-white rounded-md text-lg font-semibold
            ${loading
                            ? "bg-orange-300 cursor-not-allowed"
                            : "bg-orange-500 hover:bg-orange-600"
                        }`}
                    onClick={handleVerify}
                    disabled={loading}
                >
                    {loading ? "処理中..." : "本登録を確定する"}
                </button>

                {/* 結果メッセージ */}
                {message && (
                    <div className="mt-4 text-center text-gray-800">
                        {message}
                    </div>
                )}
            </div>
        </div>
    );
}
