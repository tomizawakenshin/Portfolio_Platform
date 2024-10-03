'use client';

import { useState } from "react";
import Header from "./components/Header";
import MainContent from "./components/MainContent";
import SignUpModal from "./components/SignUpModal";
import LoginModal from "./components/LoginModal";

export default function Home() {
  const [isModalOpen, setIsModalOpen] = useState(false);       // サインアップモーダル
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false); // ログインモーダル

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
      />

      {/* Login Modal */}
      <LoginModal
        isOpen={isLoginModalOpen}
        onClose={closeLoginModal}
        onSignUpClick={handleSwitchToSignUp}
      />
    </div>
  );
}
