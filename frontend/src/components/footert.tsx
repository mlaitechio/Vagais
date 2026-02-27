import React from "react";
// import { useNavigate } from "react-router-dom";
import section from "../assets/ailogo.png";
import linkdin from "../assets/linkdin.png";
const Footer: React.FC = () => {
  // const navigate = useNavigate();
  return (
    <footer className="footer" style={{ position: "relative", zIndex: 10, background: "linear-gradient(135deg, #0a0e27 0%, #111533 100%)" }}>
      <div className="grid_box">
        
        {/* Logo Section */}
        <div className="box">
          <a href="/" className="logo">
            <img src={section} alt="Logo" />
            <span>Powering Fintech & Financial Services</span>
          </a>
          {/* <p>
            Deploy, execute, and scale powerful AI agents in seconds.
            No coding required.
          </p> */}
          <div className="social_icon">
            <a href="#"></a>
            <a href="#"></a>
            <a href="#"></a>
            <a href="#"></a>
          </div>
        </div>

        {/* Solutions */}
        <div className="box">
          <h4>Solutions</h4>
          <ul>
            <li><a href="#">Banking</a></li>
            <li><a href="#">Fintech</a></li>
            <li><a href="#">Financial Services</a></li>
            <li><a href="#">Risk & Compliance</a></li>
          </ul>
        </div>

        {/* Company */}
        <div className="box">
          <h4>Company</h4>
          <ul>
            <li><a href="#" style={{ cursor: 'not-allowed', textDecoration: 'none' }}>Join Community</a></li>
            <li><a href="/about-us">About Us</a></li>
            <li><a href="#">Contact Us</a></li>
          </ul>
        </div>

        {/* Support */}
        <div className="box">
          <div className="support">
            <a href="tel:+919296722898">
              <span>Contact No.</span>
              +1 92967 22898
            </a>
            <a href="mailto:Agent-e@merv.one">
              <span>Email Id.</span>
              Agent-e@merv.one
            </a>
            <a href="https://www.linkedin.com/company/mlai/" className="linkdin">
              <img src={linkdin} alt="Linkdin" / >
            </a>
          </div>
        </div>

      </div>

      {/* Bottom Copy Section */}
      <div className="copywrite">
        <p> 2026 All rights reserved. <i>Product By <a href="#">MLAI Digital</a></i></p>
        <p>© 2026 All rights reserved. <i>Product By <a href="#">MLAI Digital</a></i></p>
        <ul>
          <li><a href="#">Privacy Policy</a></li>
          <li><a href="#">Terms of Service</a></li>
        </ul>
      </div>
    </footer>
  );
};

export default Footer;