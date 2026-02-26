import React from "react";
import section from "../assets/ailogo.png";
import linkdin from "../assets/linkdin.png";
const Footer: React.FC = () => {
  return (
    <footer className="footer">
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
            <li><a href="#">Join Community</a></li>
            <li><a href="#">About Us</a></li>
            <li><a href="#">Contact Us</a></li>
          </ul>
        </div>

        {/* Support */}
        <div className="box">
          <div className="support">
            <a href="tel:+91 9355641447">
              <span>Contact No.</span>
              +91 1234 567 890
            </a>
            <a href="mailto:sales@mlaitech.io">
              <span>Email Id.</span>
              info@example.com
            </a>
            <a href="https://www.linkedin.com/company/mlai/" className="linkdin">
              <img src={linkdin} alt="Linkdin" / >
            </a>
          </div>
        </div>

      </div>

      {/* Bottom Copy Section */}
      <div className="copywrite">
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