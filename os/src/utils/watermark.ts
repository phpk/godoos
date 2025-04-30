export function addWatermark(watermarkText: string = "GodoOS",container: any) {
    const watermarkDiv = document.createElement('div');
    watermarkDiv.className = 'watermark';
  
    const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    svg.setAttribute("xmlns", "http://www.w3.org/2000/svg");
    svg.setAttribute("width", "200");
    svg.setAttribute("height", "200");
    svg.setAttribute("viewBox", "0 0 200 200");
  
    const text = document.createElementNS("http://www.w3.org/2000/svg", "text");
    text.setAttribute("x", "50%");
    text.setAttribute("y", "50%");
    text.setAttribute("font-family", "Arial");
    text.setAttribute("font-size", "12");
    text.setAttribute("fill", "rgba(192, 192, 192, 0.2)");
    text.setAttribute("text-anchor", "middle");
    text.setAttribute("dominant-baseline", "middle");
    text.setAttribute("transform", "rotate(-45, 100, 100)");
    text.textContent = watermarkText;
  
    svg.appendChild(text);
    watermarkDiv.style.backgroundImage = `url('data:image/svg+xml;utf8,${encodeURIComponent(svg.outerHTML)}')`;
    watermarkDiv.style.position = 'absolute';
    watermarkDiv.style.top = '0';
    watermarkDiv.style.left = '0';
    watermarkDiv.style.width = '100%';
    watermarkDiv.style.height = '100%';
    watermarkDiv.style.zIndex = '998';
    watermarkDiv.style.pointerEvents = 'none';
    if (container.value) {
      container.value.appendChild(watermarkDiv);
    }
  };