"use client";

import React, { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { User } from "../types/User";
import { Portfolio } from "../types/Portfolio";
import Header from "../components/Header_Home";
import SkillEditModal from "../components/SkillEditModal";
import ProfileEditModal from "../components/ProfileEditModal";
import { BACKEND_URL } from "@/config";

export default function AccountPage() {
    const router = useRouter();

    // -----------------------------
    // ユーザー情報 & ポートフォリオ
    // -----------------------------
    const [user, setUser] = useState<User | null>(null);
    const [userPortfolio, setUserPortfolio] = useState<Portfolio[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    // -----------------------------
    // 学歴・プロフィール編集
    // -----------------------------
    const [isEditingAcademics, setIsEditingAcademics] = useState(false);
    const [tempSchoolName, setTempSchoolName] = useState("");
    const [tempDepartment, setTempDepartment] = useState("");
    const [tempLaboratory, setTempLaboratory] = useState("");
    const [tempGraduationYear, setTempGraduationYear] = useState("");

    // -----------------------------
    // 希望条件（複数職種）編集
    // -----------------------------
    const [isEditingDesiredJobs, setIsEditingDesiredJobs] = useState(false);
    const [tempDesiredJobTypes, setTempDesiredJobTypes] = useState<string[]>([]);
    const [jobTypeOptions, setJobTypeOptions] = useState<string[]>([]);

    // -----------------------------
    // スキル編集モーダル
    // -----------------------------
    const [isSkillModalOpen, setIsSkillModalOpen] = useState(false);

    // -----------------------------
    // 自己紹介（SelfIntroduction）編集
    // -----------------------------
    const [isEditingSelfIntro, setIsEditingSelfIntro] = useState(false);
    const [tempSelfIntro, setTempSelfIntro] = useState("");

    // -----------------------------
    // プロフィール編集モーダル（氏名・アイコン）
    // -----------------------------
    const [isProfileModalOpen, setIsProfileModalOpen] = useState(false);

    // ============================================================
    // 初期データ取得
    // ============================================================
    useEffect(() => {
        // (1) ユーザー情報
        fetch(`${BACKEND_URL}/user/GetInfo`, { credentials: "include" })
            .then((res) => {
                if (!res.ok) throw new Error("Failed to fetch user info");
                return res.json();
            })
            .then((data) => {
                setUser(data.user);
            })
            .catch((err) => console.error(err))
            .finally(() => {
                setIsLoading(false);
            });

        // (2) ポートフォリオ一覧
        fetch(`${BACKEND_URL}/Portfolio/getUserPosts`, {
            credentials: "include",
        })
            .then((res) => {
                if (!res.ok) throw new Error("Failed to fetch user portfolio");
                return res.json();
            })
            .then((data) => {
                setUserPortfolio(data.posts || []);
            })
            .catch((err) => console.error(err));

        // (3) 希望条件の職種候補
        fetch(`${BACKEND_URL}/options/job-types`, { credentials: "include" })
            .then((res) => res.json())
            .then((data) => {
                setJobTypeOptions(data.jobTypes || []);
            })
            .catch((err) => console.error(err));
    }, []);

    // ============================================================
    // 学歴・プロフィール編集ロジック
    // ============================================================
    const handleEditAcademics = () => {
        if (!user) return;
        setIsEditingAcademics(true);
        setTempSchoolName(user.SchoolName || "");
        setTempDepartment(user.Department || "");
        setTempLaboratory(user.Laboratory || "");
        setTempGraduationYear(`${user.GraduationYear}卒`);
    };

    const handleCancelAcademics = () => {
        setIsEditingAcademics(false);
    };

    const handleUpdateAcademics = async () => {
        if (!user) return;
        try {
            const numericYear = tempGraduationYear.replace("卒", "");

            const res = await fetch(`${BACKEND_URL}/user/UpdateMinimumUserInfo`, {
                method: "PUT",
                credentials: "include",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    SchoolName: tempSchoolName,
                    Department: tempDepartment,
                    Laboratory: tempLaboratory,
                    GraduationYear: numericYear,
                }),
            });
            if (!res.ok) throw new Error("学歴更新に失敗しました");

            const data = await res.json();
            setUser(data.user);
            setIsEditingAcademics(false);
        } catch (error) {
            console.error(error);
            alert("学歴の更新中にエラーが発生しました。");
        }
    };

    // ============================================================
    // 希望条件（複数職種）編集ロジック
    // ============================================================
    const handleEditDesiredJobs = () => {
        if (!user) return;
        setIsEditingDesiredJobs(true);
        setTempDesiredJobTypes(user.DesiredJobTypes || []);
    };

    const handleCancelDesiredJobs = () => {
        setIsEditingDesiredJobs(false);
    };

    const handleToggleJobType = (jobType: string) => {
        if (tempDesiredJobTypes.includes(jobType)) {
            setTempDesiredJobTypes(tempDesiredJobTypes.filter((jt) => jt !== jobType));
        } else {
            setTempDesiredJobTypes([...tempDesiredJobTypes, jobType]);
        }
    };

    const handleUpdateDesiredJobs = async () => {
        if (!user) return;
        try {
            const res = await fetch(`${BACKEND_URL}/user/UpdateMinimumUserInfo`, {
                method: "PUT",
                credentials: "include",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    DesiredJobTypes: tempDesiredJobTypes,
                }),
            });
            if (!res.ok) throw new Error("希望条件の更新に失敗しました");

            const data = await res.json();
            setUser(data.user);
            setIsEditingDesiredJobs(false);
        } catch (error) {
            console.error(error);
            alert("希望条件の更新中にエラーが発生しました。");
        }
    };

    // ============================================================
    // スキル編集モーダル
    // ============================================================
    const handleEditSkills = () => {
        setIsSkillModalOpen(true);
    };

    const handleSkillModalSave = async (newSkills: string[]) => {
        if (!user) return;
        try {
            const res = await fetch(`${BACKEND_URL}/user/UpdateMinimumUserInfo`, {
                method: "PUT",
                credentials: "include",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    Skills: newSkills,
                }),
            });
            if (!res.ok) throw new Error("スキル更新に失敗しました");

            const data = await res.json();
            setUser(data.user);
        } catch (error) {
            console.error(error);
            alert("スキルの更新中にエラーが発生しました。");
        }
    };

    // ============================================================
    // 自己紹介 (SelfIntroduction) 編集ロジック
    // ============================================================
    const handleEditIntro = () => {
        if (!user) return;
        setIsEditingSelfIntro(true);
        setTempSelfIntro(user.SelfIntroduction || "");
    };

    const handleCancelIntro = () => {
        setIsEditingSelfIntro(false);
    };

    const handleUpdateIntro = async () => {
        if (!user) return;
        try {
            const res = await fetch(`${BACKEND_URL}/user/UpdateMinimumUserInfo`, {
                method: "PUT",
                credentials: "include",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    SelfIntroduction: tempSelfIntro,
                }),
            });
            if (!res.ok) throw new Error("自己紹介の更新に失敗しました");

            const data = await res.json();
            setUser(data.user);
            setIsEditingSelfIntro(false);
        } catch (error) {
            console.error(error);
            alert("自己紹介の更新中にエラーが発生しました。");
        }
    };

    // ============================================================
    // プロフィール編集モーダル（氏名・アイコン）
    // ============================================================

    const handleEditProfileTop = () => {
        setIsProfileModalOpen(true);
    };

    // モーダルで保存成功 → user 更新
    const handleProfileSaveSuccess = (updatedUser: User) => {
        setUser(updatedUser);
    };

    // ============================================================
    // その他ハンドラ
    // ============================================================
    const handleAddWork = () => {
        router.push("/post");
    };

    const handlePortfolioClick = (postId: number) => {
        router.push(`/Portfolio/${postId}`);
    };

    // ============================================================
    // ユーティリティ
    // ============================================================
    const getProfilePhotoURL = (u: User) =>
        u.ProfileImageURL
            ? `${BACKEND_URL}/${u.ProfileImageURL}`
            : "/images/defaultUserIcon.png";

    if (isLoading) {
        return <div className="text-center mt-16">読み込み中...</div>;
    }
    if (!user) {
        return (
            <div className="text-center mt-16 text-red-500">
                ユーザー情報を取得できませんでした。
            </div>
        );
    }

    return (
        <>
            {/* ヘッダー */}
            <Header
                onPostClick={() => router.push("/post")}
                userHasPhoto={!!user.ProfileImageURL}
                userPhotoURL={getProfilePhotoURL(user)}
            />

            {/* ヘッダーの高さ分余白 */}
            <div className="pt-16" />

            {/* カバー画像 */}
            <div className="relative w-full h-40 bg-gray-200" />

            {/* ユーザーアイコン＋名前 */}
            <div className="max-w-7xl mx-auto px-4 relative">
                <div
                    className="
            absolute -top-12 left-4 
            w-24 h-24 rounded-full
            overflow-hidden border-4 border-white bg-gray-300
          "
                >
                    <img
                        src={getProfilePhotoURL(user)}
                        alt="ユーザーアイコン"
                        className="object-cover w-full h-full"
                    />
                </div>

                <div className="pt-16 flex items-center justify-between">
                    <div className="ml-32">
                        <h1 className="text-2xl font-bold">
                            {user.LastName} {user.FirstName}
                        </h1>
                        <p className="text-gray-700">
                            {user.SchoolName}・{user.GraduationYear}年卒
                        </p>
                    </div>
                    <button
                        onClick={handleEditProfileTop}
                        className="px-4 py-2 rounded bg-gray-100 hover:bg-gray-200 text-sm"
                    >
                        編集
                    </button>
                </div>

                <hr className="my-6" />
            </div>

            {/* メイン: 2カラムレイアウト */}
            <main className="max-w-7xl mx-auto px-4">
                <div className="flex flex-col md:flex-row md:space-x-8">
                    {/* 左カラム */}
                    <div className="md:w-1/3 w-full mb-8 md:mb-0 space-y-8">
                        {/* 学歴・プロフィール */}
                        <section>
                            <div className="flex items-center justify-between mb-2">
                                <h2 className="text-lg font-semibold">学歴・プロフィール</h2>
                                {!isEditingAcademics && (
                                    <button
                                        onClick={handleEditAcademics}
                                        className="px-3 py-1 rounded bg-gray-100 hover:bg-gray-200 text-xs"
                                    >
                                        編集する
                                    </button>
                                )}
                            </div>

                            {!isEditingAcademics ? (
                                // 表示モード
                                <div className="space-y-1 text-gray-800">
                                    <p>学校：{user.SchoolName}</p>
                                    <p>学部・学科：{user.Department}</p>
                                    <p>研究室：{user.Laboratory || "なし"}</p>
                                    <p>卒業予定年：{user.GraduationYear}卒</p>
                                </div>
                            ) : (
                                // 編集フォーム
                                <div className="p-4 border rounded bg-gray-50 space-y-4">
                                    <div>
                                        <label className="block text-sm font-semibold mb-1">学校</label>
                                        <input
                                            type="text"
                                            value={tempSchoolName}
                                            onChange={(e) => setTempSchoolName(e.target.value)}
                                            className="w-full p-2 border rounded"
                                        />
                                    </div>

                                    <div>
                                        <label className="block text-sm font-semibold mb-1">
                                            学部・学科・コース
                                        </label>
                                        <input
                                            type="text"
                                            value={tempDepartment}
                                            onChange={(e) => setTempDepartment(e.target.value)}
                                            className="w-full p-2 border rounded"
                                        />
                                    </div>

                                    <div>
                                        <label className="block text-sm font-semibold mb-1">
                                            研究室
                                        </label>
                                        <input
                                            type="text"
                                            value={tempLaboratory}
                                            onChange={(e) => setTempLaboratory(e.target.value)}
                                            className="w-full p-2 border rounded"
                                            placeholder="なし"
                                        />
                                    </div>

                                    <div>
                                        <label className="block text-sm font-semibold mb-1">
                                            卒業予定年
                                        </label>
                                        <select
                                            value={tempGraduationYear}
                                            onChange={(e) => setTempGraduationYear(e.target.value)}
                                            className="border rounded p-2 w-full"
                                        >
                                            <option value="2025卒">2025卒</option>
                                            <option value="2026卒">2026卒</option>
                                            <option value="2027卒">2027卒</option>
                                            <option value="2028卒">2028卒</option>
                                        </select>
                                    </div>

                                    <div className="flex justify-end space-x-2 mt-4">
                                        <button
                                            onClick={handleCancelAcademics}
                                            className="px-4 py-1 text-sm rounded bg-gray-100 hover:bg-gray-200"
                                        >
                                            キャンセル
                                        </button>
                                        <button
                                            onClick={handleUpdateAcademics}
                                            className="px-4 py-1 text-sm rounded bg-orange-500 text-white hover:bg-orange-600"
                                        >
                                            更新
                                        </button>
                                    </div>
                                </div>
                            )}
                        </section>

                        {/* 希望条件 */}
                        <section>
                            <div className="flex items-center justify-between mb-2">
                                <h2 className="text-lg font-semibold">希望条件</h2>
                                {!isEditingDesiredJobs && (
                                    <button
                                        onClick={handleEditDesiredJobs}
                                        className="px-3 py-1 rounded bg-gray-100 hover:bg-gray-200 text-xs"
                                    >
                                        編集する
                                    </button>
                                )}
                            </div>

                            {!isEditingDesiredJobs ? (
                                <div className="space-y-1 text-gray-800">
                                    {user.DesiredJobTypes?.length ? (
                                        user.DesiredJobTypes.map((jobType, idx) => (
                                            <p key={idx}>{jobType}</p>
                                        ))
                                    ) : (
                                        <p>特になし</p>
                                    )}
                                </div>
                            ) : (
                                <div className="p-4 border rounded bg-gray-50 space-y-4">
                                    <p className="text-sm text-gray-600">希望職種を複数選択できます</p>
                                    <div className="max-h-40 overflow-y-auto border p-2 rounded">
                                        {jobTypeOptions.map((option) => (
                                            <label
                                                key={option}
                                                className="flex items-center cursor-pointer mb-2"
                                            >
                                                <input
                                                    type="checkbox"
                                                    className="mr-2"
                                                    checked={tempDesiredJobTypes.includes(option)}
                                                    onChange={() => handleToggleJobType(option)}
                                                />
                                                {option}
                                            </label>
                                        ))}
                                    </div>

                                    <div className="flex justify-end space-x-2 mt-4">
                                        <button
                                            onClick={handleCancelDesiredJobs}
                                            className="px-4 py-1 text-sm rounded bg-gray-100 hover:bg-gray-200"
                                        >
                                            キャンセル
                                        </button>
                                        <button
                                            onClick={handleUpdateDesiredJobs}
                                            className="px-4 py-1 text-sm rounded bg-orange-500 text-white hover:bg-orange-600"
                                        >
                                            更新
                                        </button>
                                    </div>
                                </div>
                            )}
                        </section>

                        {/* スキル */}
                        <section>
                            <h2 className="text-lg font-semibold mb-2">スキル</h2>
                            {user.Skills?.length ? (
                                <div className="flex flex-wrap gap-2 mb-4">
                                    {user.Skills.map((skill, idx) => (
                                        <span
                                            key={idx}
                                            className="bg-gray-100 px-3 py-1 rounded-full text-sm"
                                        >
                                            {skill}
                                        </span>
                                    ))}
                                </div>
                            ) : (
                                <p className="text-gray-800 mb-4">スキルなし</p>
                            )}
                            <button
                                onClick={handleEditSkills}
                                className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 text-sm"
                            >
                                + スキルを編集
                            </button>
                        </section>
                    </div>

                    {/* 右カラム */}
                    <div className="md:w-2/3 w-full space-y-8">
                        {/* 自己紹介 */}
                        <section>
                            <div className="flex items-center justify-between mb-2">
                                <div>
                                    <h2 className="text-xl font-semibold">自己紹介</h2>
                                    <p className="text-sm text-gray-500">
                                        自己紹介は企業があなたに注目するきっかけになります！
                                    </p>
                                </div>
                                {!isEditingSelfIntro && (
                                    <button
                                        onClick={handleEditIntro}
                                        className="px-4 py-2 rounded bg-gray-100 hover:bg-gray-200 text-sm"
                                    >
                                        編集する
                                    </button>
                                )}
                            </div>

                            {!isEditingSelfIntro ? (
                                <div className="mt-2 text-gray-800">
                                    {user.SelfIntroduction || "（自己紹介がありません）"}
                                </div>
                            ) : (
                                <div className="border rounded p-4 bg-gray-50">
                                    <div className="flex justify-between mb-2">
                                        <h3 className="font-semibold">自己紹介のポイント</h3>
                                    </div>
                                    <textarea
                                        className="w-full h-32 p-2 border rounded focus:outline-none focus:ring"
                                        value={tempSelfIntro}
                                        onChange={(e) => setTempSelfIntro(e.target.value)}
                                    />
                                    <p className="text-right text-xs text-gray-500">
                                        残り{2000 - tempSelfIntro.length}文字
                                    </p>

                                    <div className="flex justify-end space-x-2 mt-3">
                                        <button
                                            onClick={handleCancelIntro}
                                            className="px-4 py-1 rounded bg-gray-100 hover:bg-gray-200 text-sm"
                                        >
                                            キャンセル
                                        </button>
                                        <button
                                            onClick={handleUpdateIntro}
                                            className="px-4 py-1 rounded bg-orange-500 text-white text-sm hover:bg-orange-600"
                                        >
                                            更新
                                        </button>
                                    </div>
                                </div>
                            )}
                        </section>

                        {/* ポートフォリオ */}
                        <section>
                            <h2 className="text-xl font-semibold mb-2">ポートフォリオ・作品</h2>
                            <p className="text-sm text-gray-500 mb-4">
                                ドラッグ＆ドロップで作品の並び替えができます。クリックで作品詳細に移動します。
                            </p>

                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
                                {userPortfolio.map((post) => {
                                    const firstImage = post.Images?.[0]?.URL
                                        ? `${BACKEND_URL}/${post.Images[0].URL}`
                                        : null;

                                    return (
                                        <div
                                            key={post.ID}
                                            className="bg-white shadow-md rounded p-4 cursor-pointer 
                                 transition-shadow duration-200 hover:shadow-lg"
                                            onClick={() => handlePortfolioClick(post.ID)}
                                        >
                                            <div className="relative overflow-hidden rounded group">
                                                {firstImage ? (
                                                    <img
                                                        src={firstImage}
                                                        alt={post.Title}
                                                        className="w-full h-40 object-cover transition-transform duration-200 group-hover:scale-105"
                                                    />
                                                ) : (
                                                    <div className="w-full h-40 bg-gray-200 flex items-center justify-center text-gray-500 rounded">
                                                        No Image
                                                    </div>
                                                )}
                                                <button
                                                    className="absolute top-2 right-2 bg-white p-1 rounded-full"
                                                    onClick={(e) => {
                                                        e.stopPropagation();
                                                        console.log("お気に入りボタンの処理");
                                                    }}
                                                >
                                                    ❤️
                                                </button>
                                            </div>
                                            <h3 className="text-lg font-semibold mt-2">{post.Title}</h3>
                                        </div>
                                    );
                                })}
                            </div>

                            <button
                                onClick={handleAddWork}
                                className="flex items-center space-x-2 px-3 py-2 
                           rounded font-bold text-white bg-orange-500 
                           hover:bg-orange-600 transition-colors duration-200"
                            >
                                <span>+ 作品を追加</span>
                            </button>
                        </section>
                    </div>
                </div>
            </main>

            {/* スキル編集モーダル */}
            {isSkillModalOpen && user && (
                <SkillEditModal
                    isOpen={isSkillModalOpen}
                    onClose={() => setIsSkillModalOpen(false)}
                    currentSkills={user.Skills || []}
                    onSave={handleSkillModalSave}
                />
            )}

            {/* プロフィール(名前/アイコン)編集モーダル */}
            {isProfileModalOpen && user && (
                <ProfileEditModal
                    isOpen={isProfileModalOpen}
                    onClose={() => setIsProfileModalOpen(false)}
                    user={user}
                    onSaveSuccess={(updatedUser) => {
                        setUser(updatedUser);
                    }}
                />
            )}
        </>
    );
}
