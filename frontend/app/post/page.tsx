"use client";

import React, { useState, useEffect, useRef } from "react";
import { useRouter } from "next/navigation";

const PostPage = () => {
    const router = useRouter();

    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");

    // スキル関連
    const [skills, setSkills] = useState<string[]>([]);
    const [skillInput, setSkillInput] = useState("");
    const [suggestedSkills, setSuggestedSkills] = useState<string[]>([]);
    const [availableSkills, setAvailableSkills] = useState<string[]>([]);
    const skillContainerRef = useRef<HTMLDivElement>(null);

    // ジャンル（複数選択ドロップダウン）
    const [genreOptions, setGenreOptions] = useState<string[]>([]);
    const [selectedGenres, setSelectedGenres] = useState<string[]>([]);
    const [isGenreDropdownOpen, setIsGenreDropdownOpen] = useState(false);
    const genreDropdownRef = useRef<HTMLDivElement>(null);

    // 画像アップロード関連
    const [images, setImages] = useState<File[]>([]);
    const fileInputRef = useRef<HTMLInputElement>(null);
    const [isDragOver, setIsDragOver] = useState(false);

    const [dirty, setDirty] = useState(false); // フォームが変更されたかどうか

    // バリデーション用エラー
    const [errors, setErrors] = useState<{
        title?: string;
        genres?: string;
        images?: string;
        skills?: string;
    }>({});

    //=================================
    // 初期処理：スキル・ジャンルリストを取得
    //=================================
    useEffect(() => {
        fetch("http://localhost:8080/options/skills", {
            credentials: "include",
        })
            .then((res) => res.json())
            .then((data) => {
                setAvailableSkills(data.skills || []);
            })
            .catch((err) => console.error("Error fetching skills:", err));

        fetch("http://localhost:8080/options/genre", {
            credentials: "include",
        })
            .then((res) => res.json())
            .then((data) => {
                setGenreOptions(data.genres || []);
            })
            .catch((err) => console.error("Error fetching genres:", err));
    }, []);

    //=================================
    // スキル候補の更新
    //=================================
    useEffect(() => {
        if (skillInput) {
            const filtered = availableSkills.filter(
                (skill) =>
                    skill.toLowerCase().includes(skillInput.toLowerCase()) &&
                    !skills.includes(skill)
            );
            setSuggestedSkills(filtered);
        } else {
            setSuggestedSkills([]);
        }
    }, [skillInput, skills, availableSkills]);

    //=================================
    // スキル候補枠を外部クリックで閉じる
    //=================================
    useEffect(() => {
        function handleSkillClickOutside(event: MouseEvent) {
            if (
                skillContainerRef.current &&
                !skillContainerRef.current.contains(event.target as Node)
            ) {
                setSuggestedSkills([]);
                setSkillInput("");
            }
        }
        document.addEventListener("mousedown", handleSkillClickOutside);
        return () => {
            document.removeEventListener("mousedown", handleSkillClickOutside);
        };
    }, []);

    //=================================
    // ジャンルのドロップダウンを外部クリックで閉じる
    //=================================
    useEffect(() => {
        function handleClickOutside(event: MouseEvent) {
            if (
                genreDropdownRef.current &&
                !genreDropdownRef.current.contains(event.target as Node)
            ) {
                setIsGenreDropdownOpen(false);
            }
        }
        document.addEventListener("mousedown", handleClickOutside);
        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, []);

    // タブを閉じる / リロード / URL直打ち(同じタブ) の場合の警告
    useEffect(() => {
        const handleBeforeUnload = (e: BeforeUnloadEvent) => {
            if (dirty) {
                // ブラウザによっては、returnValue をセットすると警告が出る
                e.preventDefault();
                e.returnValue = "";
            }
        };
        window.addEventListener("beforeunload", handleBeforeUnload);

        return () => {
            window.removeEventListener("beforeunload", handleBeforeUnload);
        };
    }, [dirty]);
    //=================================
    // スキル選択関連
    //=================================
    const handleSkillSelect = (skill: string) => {
        if (!skills.includes(skill)) {
            setSkills([...skills, skill]);
        }
        setSkillInput("");
        setSuggestedSkills([]);
    };

    const removeSkill = (skillToRemove: string) => {
        setSkills(skills.filter((skill) => skill !== skillToRemove));
    };

    const handleSkillInputKeyDown = (
        e: React.KeyboardEvent<HTMLInputElement>
    ) => {
        if (e.key === "Enter") {
            e.preventDefault();
            if (skillInput.trim()) {
                handleSkillSelect(skillInput.trim());
            }
        }
    };

    //=================================
    // ジャンル選択関連
    //=================================
    const toggleGenreDropdown = () => {
        setIsGenreDropdownOpen(!isGenreDropdownOpen);
    };

    const handleGenreChange = (genre: string) => {
        if (selectedGenres.includes(genre)) {
            setSelectedGenres(selectedGenres.filter((g) => g !== genre));
        } else {
            setSelectedGenres([...selectedGenres, genre]);
        }
    };

    //=================================
    // 画像アップロード関連
    //=================================
    const handleFiles = (files: FileList | null) => {
        if (!files) return;
        const newFiles = Array.from(files);

        if (images.length + newFiles.length > 4) {
            alert("画像は最大4枚までです。");
            return;
        }

        // 各ファイルのサイズ（最大8MB）をチェック
        const validFiles = newFiles.filter((file) => {
            if (file.size > 8 * 1024 * 1024) {
                alert(`"${file.name}" は8MBを超えています。`);
                return false;
            }
            return true;
        });

        setImages([...images, ...validFiles]);
    };

    const removeImage = (index: number) => {
        setImages(images.filter((_, i) => i !== index));
    };

    //=================================
    // 投稿ボタンクリック → バリデーション → 送信処理
    //=================================
    const handleSubmit = async () => {
        // バリデーション
        const newErrors: {
            title?: string;
            genres?: string;
            images?: string;
            skills?: string;
        } = {};

        if (!title.trim()) {
            newErrors.title = "必須項目です";
        }
        if (selectedGenres.length === 0) {
            newErrors.genres = "必須項目です";
        }
        if (images.length === 0) {
            newErrors.images = "必須項目です";
        }
        if (skills.length === 0) {
            newErrors.skills = "必須項目です";
        }

        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors);
            return;
        }

        try {
            const formData = new FormData();
            formData.append("title", title);
            formData.append("description", description);
            // ジャンル
            selectedGenres.forEach((genre) => formData.append("genres", genre));
            // スキル
            skills.forEach((skill) => formData.append("skills", skill));
            // 画像
            images.forEach((image) => formData.append("images", image, image.name));

            const response = await fetch("http://localhost:8080/Portfolio/posts", {
                method: "POST",
                credentials: "include",
                body: formData,
            });

            if (!response.ok) {
                throw new Error("投稿に失敗しました");
            }

            setDirty(false);
            // 成功時
            router.push("/home");
        } catch (error) {
            console.error("Error:", error);
        }
    };

    return (
        <div className="min-h-screen bg-white relative">
            {/* タイトル部分 */}
            <div className="max-w-5xl mx-auto px-8 pt-8 pb-4">
                <h1 className="text-2xl font-semibold mb-4">作品の投稿</h1>

                {/* 2カラムレイアウト */}
                <div className="flex flex-wrap md:flex-nowrap gap-8">
                    {/* 左カラム（画像アップ） */}
                    <div className="flex-1 min-w-[300px]">
                        <div
                            className={`border-dashed border-2 border-gray-300 p-8 flex flex-col items-center mb-2
              ${isDragOver ? "bg-gray-200" : ""}
            `}
                            onDragOver={(e) => {
                                e.preventDefault();
                                setIsDragOver(true);
                            }}
                            onDragLeave={(e) => {
                                e.preventDefault();
                                setIsDragOver(false);
                            }}
                            onDrop={(e) => {
                                e.preventDefault();
                                setIsDragOver(false);
                                handleFiles(e.dataTransfer.files);
                            }}
                            onClick={() => {
                                fileInputRef.current?.click();
                            }}
                        >
                            <p className="mb-4">クリックまたはドラッグ＆ドロップ</p>
                            <p className="text-sm text-gray-500 mb-2">
                                JPEG・PNG・GIF・PDF形式（1画像8MBまで）
                            </p>
                            <span className="text-gray-500 mb-2">残り{4 - images.length}枚</span>
                            <div className="flex space-x-4 mt-4">
                                <div className="text-center">
                                    <p>画像</p>
                                    <p className="text-xs text-gray-500">JPEG・PNG・GIF・PDF</p>
                                </div>
                            </div>
                            <input
                                type="file"
                                accept="image/*"
                                multiple
                                style={{ display: "none" }}
                                ref={fileInputRef}
                                onChange={(e) => {
                                    handleFiles(e.target.files);
                                    if (e.target) {
                                        (e.target as HTMLInputElement).value = "";
                                    }
                                }}
                            />
                        </div>
                        {/* エラーメッセージ (画像) */}
                        {errors.images && (
                            <p className="text-red-500 text-sm mb-2">{errors.images}</p>
                        )}

                        {/* 選択された画像のプレビュー */}
                        {images.length > 0 && (
                            <div className="flex flex-wrap gap-4 mb-6">
                                {images.map((image, index) => (
                                    <div key={index} className="relative">
                                        <img
                                            src={URL.createObjectURL(image)}
                                            alt={`Selected image ${index + 1}`}
                                            className="w-32 h-32 object-cover rounded"
                                        />
                                        <button
                                            type="button"
                                            className="absolute top-1 right-1 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center"
                                            onClick={() => removeImage(index)}
                                        >
                                            &times;
                                        </button>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>

                    {/* 右カラム（入力フォーム） */}
                    <div className="flex-1 min-w-[300px]">
                        {/* タイトル */}
                        <div className="mb-4">
                            <label className="block text-gray-700 mb-1">作品名</label>
                            <input
                                type="text"
                                className="w-full px-4 py-2 border rounded focus:outline-none focus:ring focus:ring-orange-500"
                                value={title}
                                onChange={(e) => setTitle(e.target.value)}
                                placeholder="作品名を入力"
                            />
                            {errors.title && (
                                <p className="text-red-500 text-sm mt-1">{errors.title}</p>
                            )}
                        </div>

                        {/* 説明文 */}
                        <div className="mb-4">
                            <label className="block text-gray-700 mb-1">説明文</label>
                            <textarea
                                className="w-full px-4 py-2 border rounded focus:outline-none focus:ring focus:ring-orange-500"
                                rows={4}
                                value={description}
                                onChange={(e) => setDescription(e.target.value)}
                                placeholder="作品の説明をこちらにご記入ください。"
                            />
                        </div>

                        {/* ジャンル (複数選択) */}
                        <div className="mb-4 relative" ref={genreDropdownRef}>
                            <label className="block text-gray-700 mb-1">作品のジャンル</label>
                            <div
                                className="border border-gray-300 rounded p-2 relative flex items-center justify-between cursor-pointer"
                                onClick={toggleGenreDropdown}
                            >
                                <span className="text-gray-700">
                                    {selectedGenres.length > 0
                                        ? selectedGenres.join("、 ")
                                        : "ジャンルを選択"}
                                </span>
                                <span className="text-gray-600 ml-2">▼</span>
                            </div>
                            {errors.genres && (
                                <p className="text-red-500 text-sm mt-1">{errors.genres}</p>
                            )}
                            {isGenreDropdownOpen && (
                                <div className="absolute left-0 bg-white border border-gray-300 rounded mt-2 w-full max-h-40 overflow-y-auto z-50">
                                    {genreOptions.map((genre) => (
                                        <label
                                            key={genre}
                                            className="flex items-center p-2 hover:bg-gray-200 cursor-pointer"
                                        >
                                            <input
                                                type="checkbox"
                                                checked={selectedGenres.includes(genre)}
                                                onChange={() => handleGenreChange(genre)}
                                                className="mr-2"
                                            />
                                            {genre}
                                        </label>
                                    ))}
                                </div>
                            )}
                        </div>

                        {/* スキル入力 */}
                        <div className="mb-4 relative" ref={skillContainerRef}>
                            <label className="block text-gray-700 mb-1">
                                使用したツール・スキル
                            </label>
                            <div className="border border-gray-300 rounded p-2 relative">
                                {/* 選択されたスキルをタグとして表示 */}
                                <div className="flex flex-wrap gap-2 mb-2">
                                    {skills.map((skill) => (
                                        <span
                                            key={skill}
                                            className="bg-gray-200 px-3 py-1 rounded-full text-sm flex items-center"
                                        >
                                            {skill}
                                            <button
                                                className="ml-2 text-red-500"
                                                onClick={() => removeSkill(skill)}
                                            >
                                                &times;
                                            </button>
                                        </span>
                                    ))}
                                </div>
                                <input
                                    type="text"
                                    placeholder="スキルを選択・追加"
                                    value={skillInput}
                                    onChange={(e) => setSkillInput(e.target.value)}
                                    onKeyDown={handleSkillInputKeyDown}
                                    className="w-full p-2 border border-gray-300 rounded"
                                />
                                {/* スキル候補ドロップダウン */}
                                {skillInput && (
                                    <div
                                        className="absolute left-0 z-50 bg-white border border-gray-300 rounded mt-2 w-full max-h-40 overflow-y-auto"
                                        onClick={(e) => e.stopPropagation()}
                                    >
                                        {suggestedSkills.length > 0 ? (
                                            suggestedSkills.map((skill) => (
                                                <div
                                                    key={skill}
                                                    className="p-2 hover:bg-gray-200 cursor-pointer"
                                                    onClick={() => handleSkillSelect(skill)}
                                                >
                                                    {skill}
                                                </div>
                                            ))
                                        ) : (
                                            <div
                                                className="p-2 hover:bg-gray-200 cursor-pointer"
                                                onClick={() => handleSkillSelect(skillInput.trim())}
                                            >
                                                Enterで追加
                                            </div>
                                        )}
                                    </div>
                                )}
                            </div>
                            {errors.skills && (
                                <p className="text-red-500 text-sm mt-1">{errors.skills}</p>
                            )}
                        </div>
                    </div>
                </div>
            </div>

            {/* 下のボタン領域との仕切り線 */}
            <hr className="border-t border-gray-200 mt-4" />

            {/* 画面下部固定の投稿ボタン */}
            <div
                className="
          fixed
          bottom-0
          left-0
          w-full
          border-t
          border-gray-200
          bg-white
          p-2
          flex
          justify-end
        "
            >
                <button
                    onClick={handleSubmit}
                    className="
            px-4
            py-2
            text-sm
            font-bold
            text-white
            bg-orange-500
            hover:bg-orange-600
            rounded
          "
                >
                    投稿する
                </button>
            </div>
        </div>
    );
};

export default PostPage;
