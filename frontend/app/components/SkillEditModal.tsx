"use client";

import { BACKEND_URL } from "@/config";
import React, { useEffect, useState, useRef } from "react";

interface SkillEditModalProps {
    isOpen: boolean;
    onClose: () => void;
    currentSkills: string[];  // 現在のユーザーのスキル一覧
    onSave: (newSkills: string[]) => void;
    // 親コンポーネント側でスキル更新APIを叩き、Userを更新する
}

const SkillEditModal: React.FC<SkillEditModalProps> = ({
    isOpen,
    onClose,
    currentSkills,
    onSave,
}) => {
    const [skills, setSkills] = useState<string[]>([]);
    const [skillInput, setSkillInput] = useState("");
    const [suggestedSkills, setSuggestedSkills] = useState<string[]>([]);
    const [availableSkills, setAvailableSkills] = useState<string[]>([]);
    const skillContainerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (isOpen) {
            // モーダルが開くたびに currentSkills を初期化
            setSkills(currentSkills);
            setSkillInput("");
        }
    }, [isOpen, currentSkills]);

    // スキル候補をバックエンドから取得
    useEffect(() => {
        fetch(`${BACKEND_URL}/options/skills`, { credentials: "include" })
            .then((res) => res.json())
            .then((data) => {
                setAvailableSkills(data.skills || []);
            })
            .catch((err) => console.error(err));
    }, []);

    // 入力補完（suggestedSkillsの更新）
    useEffect(() => {
        if (skillInput) {
            const filtered = availableSkills.filter(
                (s) =>
                    s.toLowerCase().includes(skillInput.toLowerCase()) &&
                    !skills.includes(s)
            );
            setSuggestedSkills(filtered);
        } else {
            setSuggestedSkills([]);
        }
    }, [skillInput, skills, availableSkills]);

    // 外側クリックで閉じる or suggestionリセット
    useEffect(() => {
        function handleClickOutside(event: MouseEvent) {
            if (
                skillContainerRef.current &&
                !skillContainerRef.current.contains(event.target as Node)
            ) {
                setSuggestedSkills([]);
            }
        }
        document.addEventListener("mousedown", handleClickOutside);
        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, []);

    const handleAddSkill = (skill: string) => {
        if (!skills.includes(skill)) {
            setSkills([...skills, skill]);
        }
        setSkillInput("");
        setSuggestedSkills([]);
    };

    const removeSkill = (skill: string) => {
        setSkills(skills.filter((s) => s !== skill));
    };

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") {
            e.preventDefault();
            if (skillInput.trim()) {
                handleAddSkill(skillInput.trim());
            }
        }
    };

    const handleSave = () => {
        // 親コンポーネントに newSkills を渡す
        onSave(skills);
        onClose();
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            {/* モーダル枠 */}
            <div className="bg-white rounded p-6 w-full max-w-md relative">
                <button
                    className="absolute top-2 right-2 text-gray-500 hover:text-gray-700"
                    onClick={onClose}
                >
                    &times;
                </button>
                <h2 className="text-xl font-semibold mb-4">スキルを編集</h2>

                {/* 選択済みスキル */}
                <div className="flex flex-wrap gap-2 mb-4">
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

                {/* 入力 + サジェスト */}
                <div className="relative mb-4" ref={skillContainerRef}>
                    <input
                        type="text"
                        value={skillInput}
                        onChange={(e) => setSkillInput(e.target.value)}
                        onKeyDown={handleKeyDown}
                        placeholder="スキルを選択・追加"
                        className="w-full border p-2 rounded"
                    />
                    {skillInput && suggestedSkills.length > 0 && (
                        <div className="absolute left-0 mt-1 bg-white border rounded shadow w-full max-h-40 overflow-y-auto z-50">
                            {suggestedSkills.map((s) => (
                                <div
                                    key={s}
                                    className="px-2 py-1 hover:bg-gray-100 cursor-pointer"
                                    onClick={() => handleAddSkill(s)}
                                >
                                    {s}
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                {/* 保存ボタン */}
                <div className="flex justify-end">
                    <button
                        onClick={handleSave}
                        className="px-4 py-2 bg-orange-500 text-white rounded hover:bg-orange-600"
                    >
                        保存する
                    </button>
                </div>
            </div>
        </div>
    );
};

export default SkillEditModal;
