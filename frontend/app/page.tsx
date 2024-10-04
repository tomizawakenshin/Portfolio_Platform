'use client';

import { useState } from "react";
import Header from "./components/Header";
import MainContent from "./components/MainContent";
import SignUpModal from "./components/SignUpModal";
import LoginModal from "./components/LoginModal";
import EmailSignUpModal from "./components/EmailSignUpModal";

export default function Home() {
  const [isModalOpen, setIsModalOpen] = useState(false);       // サインアップモーダル
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false); // ログインモーダル
  const [isEmailSignUpModalOpen, setIsEmailSignUpModalOpen] = useState(false); // 新しいメールサインアップモーダル

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
      />
    </div>
  );
}
