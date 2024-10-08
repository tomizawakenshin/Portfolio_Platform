'use client';

import { useState } from "react";
import Header from "../components/Header";
import MainContent from "../components/MainContent";
import SignUpModal from "../components/SignUpModal";
import LoginModal from "../components/LoginModal";
import EmailSignUpModal from "../components/EmailSignUpModal";
import SignUpCompleteModal from "../components/SignUpCompleteModal";

export default function Home() {
  const [isModalOpen, setIsModalOpen] = useState(false);       // サインアップモーダル
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false); // ログインモーダル
  const [isEmailSignUpModalOpen, setIsEmailSignUpModalOpen] = useState(false); // 新しいメールサインアップモーダル
  const [isCompleteModalOpen, setIsCompleteModalOpen] = useState(false);
  const [userEmail, setUserEmail] = useState("");

  const openSignUpModal = () => {
    setIsModalOpen(true);
  };

  const closeSignUpModal = () => {
    setIsModalOpen(false);
  };

  const openLoginModal = () => {
    setIsLoginModalOpen(true);
  };

  const closeLoginModal = () => {
    setIsLoginModalOpen(false);
  };

  // SignUpModalでログインボタンをクリックしたときの処理
  const handleSwitchToLogin = () => {
    closeSignUpModal();       // サインアップモーダルを閉じる
    openLoginModal();   // ログインモーダルを開く
  };

  const handleSwitchToSignUp = () => {
    closeLoginModal();
    openSignUpModal();
  }

  const openEmailSignUpModal = () => {
    setIsEmailSignUpModalOpen(true);
    closeSignUpModal();
  };

  const closeEmailSignUpModal = () => {
    setIsEmailSignUpModalOpen(false);
  };

  const handleSignUpComplete = (email: string) => {
    setUserEmail(email);
    setIsCompleteModalOpen(true);
  };

  const closeCompleteModal = () => {
    setIsCompleteModalOpen(false);
  };

  return (
    <div className="font-sans">
      {/* Header */}
      <Header />

      {/* Main content */}
      <MainContent
        onFreeStartClick={openSignUpModal}
        onLoginClick={openLoginModal}
      />

      {/* SignUp Modal */}
      <SignUpModal
        isOpen={isModalOpen}
        onClose={closeSignUpModal}
        onLoginClick={handleSwitchToLogin} // ログインボタンがクリックされたときのハンドラ
        onEmailSignUpClick={openEmailSignUpModal} // メールサインアップボタンを追加
      />

      {/* Login Modal */}
      <LoginModal
        isOpen={isLoginModalOpen}
        onClose={closeLoginModal}
        onSignUpClick={handleSwitchToSignUp}
      />

      <EmailSignUpModal
        isOpen={isEmailSignUpModalOpen}
        onClose={closeEmailSignUpModal}
        onComplete={handleSignUpComplete}
      />

      <SignUpCompleteModal
        isOpen={isCompleteModalOpen}
        onClose={closeCompleteModal}
        email={userEmail}
      />
    </div>
  );
}
