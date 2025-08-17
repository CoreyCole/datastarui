var L=`<style>
  :host {
    --inspector-bg: #ebede9;
    --inspector-bg-light: #c7cfcc;
    --inspector-bg-hover: #d9dedc;
    --inspector-bg-dark: #a8b5b2;
    --inspector-color: #090a14;
    --inspector-color-blur: #819796;
    --inspector-color-filtered: #d73a49;
    --inspector-color-matched: #28a745;
    --color-purple: #A599FF;
    --color-gold: #FFA116;
    --color-html5: #E34C26;
    --color-html5-dark: #F06529;
    --color-dark-text: #0A0A0A;
    --color-white: #FFFFFF;
    background-color: var(--inspector-bg);
    border: 1px solid var(--inspector-bg-light);
    border-bottom: 0;
    border-radius: 7px 7px 0 0;
    color: var(--inspector-color);
    position: fixed;
    bottom: 0;
    right: 24px;
    max-width: 640px;
    margin: 0 5%;
    font-family: monospace;
    z-index: 9999;

    header, #collapsed {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 20px;
      padding: 8px 10px;
      font-size: 13px;
      font-weight: bold;
      user-select: none;
      cursor: pointer;
      box-sizing: border-box;
      height: 40px;

      button {
        padding: 2px 4px;
      }

      h1 {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: inherit;
        margin: 0;
        margin-left: -3px;

        button {
          font-size: 17px;
          padding: 0px 4px;
        }
      }
    }

    header {
      width: 640px;

      nav {
        display: flex;
        gap: 4px;
        
        button {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 4px;

          span {
            font-weight: bold;
          }
          
          svg {
            color: inherit;
          }
          
          &[data-type="currentSignals"][aria-selected="true"] {
            background-color: var(--color-purple);
            color: var(--color-dark-text);
            border-color: var(--color-purple);
          }
          
          &[data-type="signalPatchEvent"][aria-selected="true"] {
            background-color: var(--color-gold);
            color: var(--color-dark-text);
            border-color: var(--color-gold);
          }
          
          &[data-type="sseEvent"][aria-selected="true"] {
            background-color: var(--color-html5);
            color: var(--color-white);
            border-color: var(--color-html5);
          }
          
          &[data-type="persistedData"][aria-selected="true"] {
            background-color: var(--inspector-color-matched);
            color: var(--color-white);
            border-color: var(--inspector-color-matched);
          }
        }
      }
    }

    #collapsed {
      padding: 8px 12px 8px 15px;
    }

    button {
      background-color: var(--inspector-bg);
      border: 1px solid var(--inspector-bg-dark);
      border-radius: 4px;
      color: var(--inspector-color);
      font-family: inherit;
      font-size: 11px;
      flex-shrink: 0;
      cursor: pointer;

      &:hover {
        background-color: var(--inspector-bg-hover);
      }
    }

    pre {
      margin: 0;
    }

    main {
      max-height: 40vh;
      overflow-y: auto;
      font-size: 12px;
      border-top: 1px solid var(--inspector-bg-dark);

      .filters input[type="text"] {
        background-color: var(--inspector-bg);
        border: 1px solid var(--inspector-bg-dark);
        border-radius: 4px;
        color: var(--inspector-color);
        font-family: inherit;
        font-size: 11px;
        padding: 4px 6px;
        flex: 1;
        min-width: 0;
      }
      
      .input-with-icon {
        position: relative;
        display: flex;
        flex: 1;
        min-width: 0;
        
        input {
          padding-left: 32px;
          width: 100%;
        }
        
        .input-icon {
          position: absolute;
          left: 8px;
          top: 50%;
          transform: translateY(-50%);
          width: 16px;
          height: 16px;
          pointer-events: none;
          opacity: 0.6;
          color: var(--inspector-color);
        }
      }

      .filters label {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 11px;
        user-select: none;
        cursor: pointer;
        white-space: nowrap;
        margin: 0 4px;
        
        span {
          font-weight: bold;
          font-family: monospace;
        }
      }

      .filters input[type="checkbox"] {
        cursor: pointer;
        margin: 0;
        flex-shrink: 0;
      }
        
      #resetFilters, #resetPersistedDataFilters, #refreshPersistedData {
        padding-bottom: 0;
      }

      #copyFilterObject, #copyPersistedDataFilterObject {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 4px;
        flex-shrink: 0;
        gap: 4px;
        
        svg {
          width: 14px;
          height: 14px;
        }
      }

      #clearPersistedData {
        display: flex; 
        align-items: center; 
        gap: 2px;
      }

      .regex-display {
        display: flex;
        align-items: center;
        gap: 8px;
        width: 100%;
        font-size: 11px;
        font-family: monospace;
        
        .regex-content {
          display: flex;
          gap: 4px;
          flex: 1;
          
          &::before {
            content: "{";
            color: var(--inspector-color);
          }
          
          &::after {
            content: "}";
            color: var(--inspector-color);
          }
        }
        
        .regex-include {
          color: var(--inspector-color-matched);
        }
        
        .regex-exclude {
          color: var(--inspector-color-filtered);
        }
        
        .regex-include:empty,
        .regex-exclude:empty {
          display: none;
        }
        
        .regex-include:not(:empty)::before {
          content: "include: ";
          color: var(--inspector-color);
          font-weight: normal;
        }
        
        .regex-exclude:not(:empty)::before {
          content: ", exclude: ";
          color: var(--inspector-color);
          font-weight: normal;
        }
      }

      p {
        margin: 8px 0;
        padding: 4px 10px;
      }

      footer {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 8px 10px;
        border-top: 1px solid var(--inspector-bg-dark);
        background-color: var(--inspector-bg-hover);
        font-size: 11px;
        
        &.filters {
          flex-wrap: wrap;
          row-gap: 12px;
          column-gap: 8px;
        }
        
        label {
          font-weight: normal;
        }
      }

      ul {
        margin: 8px 0;
        padding: 0;
        list-style: none;

        li {
          display: flex;
          align-items: flex-start;
          gap: 4px;
          padding: 4px 10px;

          &:hover {
            background-color: var(--inspector-bg-hover);
          }

          span {
            flex-shrink: 0;
            margin-right: 4px;
          }

          div {
            display: flex;
            flex: 1;
            overflow-x: auto;
          }
          
          summary {
            list-style: none;
            position: relative;
            padding-right: 1.2em;
            cursor: pointer;
            width: fit-content;
            
            &::after {
              content: '\u25B8';
              position: absolute;
              top: 0;
              right: 0;
              font-size: 16px;
              line-height: 16px;
              transition: transform 0.2s;
            }
          }

          details[open] {
            summary.blurrable {
              color: var(--inspector-color-blur);
            }
            
            summary::after {
              transform: rotate(90deg);
            }
          }
        }
      }

      .currentSignals ul li {
        border-left: 3px solid var(--color-purple);

        &:hover {
          background-color: initial;
        }
      }

      .signalPatchEvents ul li {
        border-left: 3px solid var(--color-gold);
      }

      .sseEvents ul li {
        border-left: 3px solid var(--color-html5-dark);
      }

      .persistedData ul li {
        border-left: 3px solid var(--inspector-color-matched);
      }

      button {
        padding: 2px 4px;
        flex-shrink: 0;
        
        &.copy-btn, &.log-btn, &.clear-btn {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 4px;
          
          svg {
            width: 16px;
            height: 16px;
          }
        }
      }
      
      input[type="number"] {
        background-color: var(--inspector-bg);
        border: 1px solid var(--inspector-bg-dark);
        border-radius: 4px;
        color: var(--inspector-color);
        font-family: inherit;
        font-size: 11px;
        padding: 4px 6px;
        flex-shrink: 0;
      }
      
      pre {
        .filtered {
          color: var(--inspector-color-filtered);
        }
        .matched {
          color: var(--inspector-color-matched);
        }
      }
      
      .signals-table {
        width: 100%;
        border-collapse: collapse;
        font-family: monospace;
        
        th {
          background-color: var(--inspector-bg-dark);
          padding: 6px 10px;
          text-align: left;
          border: 1px solid var(--inspector-bg-light);
          font-weight: bold;
        }
        
        td {
          padding: 4px 10px;
          border: 1px solid var(--inspector-bg-light);
        }
        
        tr:nth-child(even) {
          background-color: var(--inspector-bg-hover);
        }
        
        tr.matched td {
          color: var(--inspector-color-matched);
        }
        
        tr.filtered td {
          color: var(--inspector-color-filtered);
        }
      }
      
      /* Datastar signal highlighting styles */
      .datastar-signal-highlight {
        outline: 3px solid var(--color-purple) !important;
        outline-offset: 3px !important;
        background-color: rgba(165, 153, 255, 0.3) !important;
        box-shadow: 0 0 10px rgba(165, 153, 255, 0.5) !important;
        border: 2px solid var(--color-purple) !important;
        position: relative !important;
        z-index: 9998 !important;
      }
      
      .signal-path-hoverable {
        cursor: pointer;
        transition: background-color 0.2s;
      }
      
      .signal-path-hoverable:hover {
        background-color: var(--inspector-bg-hover);
      }

      .toggle-switch {
        position: relative;
        display: inline-block;
        width: 40px;
        height: 20px;
        
        input {
          opacity: 0;
          width: 0;
          height: 0;
          
          &:checked + .toggle-slider {
            background-color: var(--inspector-color);
          }
          
          &:checked + .toggle-slider:before {
            transform: translateX(18px);
          }
        }
        
        .toggle-slider {
          position: absolute;
          cursor: pointer;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background-color: var(--inspector-bg-dark);
          transition: .2s;
          border-radius: 20px;
          
          &:before {
            position: absolute;
            content: "";
            height: 16px;
            width: 16px;
            left: 2px;
            bottom: 2px;
            background-color: var(--inspector-bg);
            transition: .2s;
            border-radius: 50%;
          }
        }
      }
      
      
      .filter-row {
        display: flex;
        align-items: center;
        gap: 8px;
        width: 100%;
      }
      
      section#persistedDataContent {
        h3 {
          display: flex; 
          align-items: center; 
          justify-content: center; 
          gap: 8px;
          margin: 8px 10px; 
          font-size: 12px; 
          font-weight: bold; 
          color: var(--inspector-color);

        }

        #persistedDataStorageTitle {
          display: flex; 
          align-items: center; 
          justify-content: center; 
        }

        .persisted-key-section {
          width: 100%;
          display: flex;
          flex-direction: column;
          border: 1px solid var(--inspector-bg-dark);
          border-radius: 4px;
        }
        
        .persisted-key-caption {
          font-weight: bold;
          padding: 4px 8px;
          background-color: var(--inspector-bg-hover);
          border-radius: 4px 4px 0 0;
          border-bottom: 1px solid var(--inspector-bg-dark);
          display: flex;
          align-items: center;
          justify-content: end;
          width: 100%;
          box-sizing: border-box;
        }
        
        .caption-text {
          font-size: 12px;
          color: var(--inspector-color);
        }
        
        .caption-buttons {
          display: flex;
          align-items: center;
          justify-content: end;
          gap: 4px;
          flex-shrink: 0;
        }
        
        pre {
          padding: 8px;
          margin: 0;
        }
        
        table {
          width: 100%;
          border-collapse: collapse;
          font-family: monospace;
        }
        
        table th {
          background-color: var(--inspector-bg-dark);
          padding: 6px 10px;
          text-align: left;
          font-weight: bold;
        }
        
        table td {
          padding: 4px 10px;
        }
        
        table tr:nth-child(even) {
          background-color: var(--inspector-bg-hover);
        }
        
        table tr.matched td {
          color: var(--inspector-color-matched);
        }
        
        table tr.filtered td {
          color: var(--inspector-color-filtered);
        }
      }
      
      .filter-row-secondary {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 8px;
        width: 100%;
      }
    }
  }


  @media (prefers-color-scheme: dark) {
    :host {
      --inspector-bg: #202e37;
      --inspector-bg-light: #394a50;
      --inspector-bg-hover: #1a262e;
      --inspector-bg-dark: #577277;
      --inspector-color: #ebede9;
      --inspector-color-blur: #819796;
      --inspector-color-filtered: #f97583;
      --inspector-color-matched: #34d058;
    }

    header nav button {
      &[data-type="currentSignals"][aria-selected="true"] {
        background-color: var(--color-purple);
        color: var(--color-dark-text);
        border-color: var(--color-purple);
      }
      
      &[data-type="signalPatchEvent"][aria-selected="true"] {
        background-color: var(--color-gold);
        color: var(--color-dark-text);
        border-color: var(--color-gold);
      }
      
      &[data-type="sseEvent"][aria-selected="true"] {
        background-color: var(--color-html5);
        color: var(--color-white);
        border-color: var(--color-html5);
      }
      
      &[data-type="persistedData"][aria-selected="true"] {
        background-color: var(--inspector-color-matched);
        color: var(--color-white);
        border-color: var(--inspector-color-matched);
      }
    }
  }
</style>

<template id="copy-button-template">
  <button class="copy-btn" title="Copy to clipboard">
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
      <use href="#copy-icon"/>
    </svg>
    <span style="display: none"></span>
  </button>
</template>

<template id="log-button-template">
  <button class="log-btn" title="Log to the console">
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
      <path fill="currentColor" fill-rule="evenodd" d="M9.586 11L7.05 8.464L8.464 7.05l4.95 4.95l-4.95 4.95l-1.414-1.414L9.586 13H3v-2zM11 3h8c1.1 0 2 .9 2 2v14c0 1.1-.9 2-2 2h-8v-2h8V5h-8z"/>
    </svg>
    <span style="display: none"></span>
  </button>
</template>

<template id="clear-button-template">
  <button class="clear-btn" title="Clear persisted data">
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
      <use href="#clear-icon"/>
    </svg>
    <span>Clear</span>
  </button>
</template>

<svg style="display: none">
  <defs>
    <g id="datastar-logo">
      <g fill="#752438"><path d="m20 5h1v2h-1z"/><path d="m19 6h1v2h-1z"/><path d="m21 3h1v2h-1z"/><path d="m5 13h1v1h-1z"/><path d="m10 19h1v1h-1z"/><path d="m14 15h1v1h-1z"/><path d="m13 16h1v1h-1z"/><path d="m12 17h1v1h-1z"/><path d="m4 14h1v1h-1z"/><path d="m17 6h1v1h-1z"/></g>
      <g fill="#73bed3"><path d="m11 7h1v1h-1z"/><path d="m12 6h1v1h-1z"/><path d="m13 5h1v1h-1z"/><path d="m9 9h1v1h-1z"/><path d="m10 8h1v1h-1z"/></g>
      <g fill="#172038"><path d="m9 12h1v1h-1z"/><path d="m16 12h1v1h-1z"/><path d="m12 15h1v2h-1z"/><path d="m15 13h1v1h-1z"/><path d="m13 15h1v1h-1z"/><path d="m8 11h1v1h-1z"/><path d="m17 11h1v1h-1z"/><path d="m10 17h2v1h-2z"/><path d="m6 14h1v1h-1z"/><path d="m11 14h1v1h-1z"/><path d="m14 14h1v1h-1z"/><path d="m18 10h1v1h-1z"/></g>
      <g fill="#be772b"><path d="m8 19h1v2h-1z"/><path d="m5 20h2v1h-2z"/><path d="m9 18h1v1h-1z"/><path d="m5 16h1v1h-1z"/></g>
      <g fill="#da863e"><path d="m17 2h3v1h-3z"/><path d="m15 4h1v1h-1z"/><path d="m16 3h1v1h-1z"/></g>
      <g fill="#cf573c"><path d="m6 10h1v1h-1z"/><path d="m4 12h1v1h-1z"/><path d="m10 20h1v1h-1z"/><path d="m3 13h1v1h-1z"/><path d="m9 13h1v1h-1z"/><path d="m20 2h2v1h-2z"/><path d="m11 19h1v1h-1z"/><path d="m6 16h1v1h-1z"/><path d="m14 16h1v1h-1z"/><path d="m7 15h1v1h-1z"/><path d="m13 17h1v1h-1z"/><path d="m18 3h1v1h-1z"/><path d="m5 11h1v1h-1z"/><path d="m2 14h1v1h-1z"/><path d="m8 14h1v1h-1z"/><path d="m17 3h1v2h-1z"/><path d="m12 18h1v1h-1z"/></g>
      <g fill="#3c5e8b"><path d="m11 9h1v3h-1z"/><path d="m7 14h1v1h-1z"/><path d="m15 11h1v1h-1z"/><path d="m11 15h1v1h-1z"/><path d="m14 10h1v3h-1z"/><path d="m12 8h1v5h-1z"/><path d="m8 12h1v2h-1z"/><path d="m13 7h1v1h-1z"/><path d="m16 7h1v4h-1z"/><path d="m6 13h1v1h-1z"/><path d="m15 6h1v2h-1z"/><path d="m13 11h1v3h-1z"/><path d="m14 6h1v3h-1z"/><path d="m13 9h1v1h-1z"/><path d="m15 9h1v1h-1z"/><path d="m17 8h1v2h-1z"/><path d="m10 15h1v2h-1z"/><path d="m8 10h1v1h-1z"/><path d="m10 10h1v1h-1z"/></g>
      <g fill="#253a5e"><path d="m17 10h1v1h-1z"/><path d="m11 13h1v1h-1z"/><path d="m14 13h1v1h-1z"/><path d="m18 9h1v1h-1z"/><path d="m9 16h1v1h-1z"/><path d="m11 16h1v1h-1z"/><path d="m7 12h1v2h-1z"/><path d="m12 14h1v1h-1z"/><path d="m9 11h1v1h-1z"/><path d="m16 11h1v1h-1z"/><path d="m10 12h1v1h-1z"/><path d="m15 12h1v1h-1z"/></g>
      <g fill="#e8c170"><path d="m8 17h1v1h-1z"/><path d="m6 15h1v1h-1z"/><path d="m4 19h1v2h-1z"/><path d="m5 18h1v2h-1z"/><path d="m3 20h1v2h-1z"/><path d="m7 18h1v1h-1z"/></g>
      <g fill="#241527"><path d="m20 7h1v1h-1z"/><path d="m18 8h1v1h-1z"/><path d="m19 8h1v2h-1z"/><path d="m21 5h1v2h-1z"/></g>
      <g fill="#602c2c"><path d="m9 17h1v1h-1z"/><path d="m6 18h1v1h-1z"/></g>
      <g fill="#4f8fba"><path d="m11 8h1v1h-1z"/><path d="m7 11h1v1h-1z"/><path d="m10 11h1v1h-1z"/><path d="m13 6h1v1h-1z"/><path d="m13 14h1v1h-1z"/><path d="m9 10h1v1h-1z"/><path d="m6 12h1v1h-1z"/><path d="m11 12h1v1h-1z"/><path d="m14 5h1v1h-1z"/><path d="m10 9h1v1h-1z"/><path d="m12 13h1v1h-1z"/><path d="m12 7h1v1h-1z"/></g>
      <g fill="#a53030"><path d="m19 3h1v3h-1z"/><path d="m6 11h1v1h-1z"/><path d="m6 17h1v1h-1z"/><path d="m14 17h1v1h-1z"/><path d="m13 18h1v1h-1z"/><path d="m20 3h1v2h-1z"/><path d="m3 14h1v1h-1z"/><path d="m9 14h1v1h-1z"/><path d="m12 19h1v1h-1z"/><path d="m8 15h1v1h-1z"/><path d="m7 16h1v1h-1z"/><path d="m4 13h1v1h-1z"/><path d="m10 13h1v1h-1z"/><path d="m18 4h1v3h-1z"/><path d="m7 10h1v1h-1z"/><path d="m16 4h1v2h-1z"/><path d="m11 20h1v1h-1z"/><path d="m17 5h1v1h-1z"/><path d="m5 12h1v1h-1z"/><path d="m10 21h1v1h-1z"/><path d="m14 4h1v1h-1z"/></g>
      <g fill="#411d31"><path d="m15 5h1v1h-1z"/><path d="m8 16h1v1h-1z"/><path d="m9 15h1v1h-1z"/><path d="m17 7h2v1h-2z"/><path d="m16 6h1v1h-1z"/><path d="m5 14h1v1h-1z"/><path d="m10 14h1v1h-1z"/><path d="m10 18h2v1h-2z"/><path d="m7 17h1v1h-1z"/></g>
      <g fill="#a4dddb"><path d="m13 10h1v1h-1z"/><path d="m15 10h1v1h-1z"/><path d="m13 8h1v1h-1z"/><path d="m15 8h1v1h-1z"/><path d="m14 9h1v1h-1z"/></g>
      <g fill="#de9e41"><path d="m4 18h1v1h-1z"/><path d="m8 18h1v1h-1z"/><path d="m5 17h1v1h-1z"/><path d="m3 19h1v1h-1z"/><path d="m6 19h2v1h-2z"/></g>
    </g>
    
    <!-- Common icons -->
    <g id="copy-icon">
      <path fill="currentColor" d="M4 2h11v2H6v13H4zm4 4h12v16H8zm2 2v12h8V8z"/>
    </g>
    
    <g id="reset-icon">
      <path fill="currentColor" d="M11.77 3c-2.65.07-5 1.28-6.6 3.16L3.85 4.85a.5.5 0 0 0-.85.36V9.5c0 .28.22.5.5.5h4.29c.45 0 .67-.54.35-.85L6.59 7.59C7.88 6.02 9.82 5 12 5c4.32 0 7.74 3.94 6.86 8.41c-.54 2.77-2.81 4.98-5.58 5.47c-3.8.68-7.18-1.74-8.05-5.16c-.12-.42-.52-.72-.96-.72c-.65 0-1.14.61-.98 1.23C4.28 18.12 7.8 21 12 21c5.06 0 9.14-4.17 9-9.26c-.14-4.88-4.35-8.86-9.23-8.74M14 12c0-1.1-.9-2-2-2s-2 .9-2 2s.9 2 2 2s2-.9 2-2" />
    </g>
    
    <g id="refresh-icon">
      <path fill="currentColor" d="M17.65 6.35C16.2 4.9 14.21 4 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08c-.82 2.33-3.04 4-5.65 4c-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4z"/>
    </g>
    
    <g id="clear-icon">
      <path fill="currentColor" d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z"/>
    </g>
    
    <g id="tree-icon">
      <path fill="currentColor" d="M240 128c-64 0 0 88-64 88H80c-64 0 0-88-64-88c64 0 0-88 64-88h96c64 0 0 88 64 88" opacity=".2"/>
      <path fill="currentColor" d="M43.18 128a29.8 29.8 0 0 1 8 10.26c4.8 9.9 4.8 22 4.8 33.74c0 24.31 1 36 24 36a8 8 0 0 1 0 16c-17.48 0-29.32-6.14-35.2-18.26c-4.8-9.9-4.8-22-4.8-33.74c0-24.31-1-36-24-36a8 8 0 0 1 0-16c23 0 24-11.69 24-36c0-11.72 0-23.84 4.8-33.74C50.68 38.14 62.52 32 80 32a8 8 0 0 1 0 16c-23 0-24 11.69-24 36c0 11.72 0 23.84-4.8 33.74A29.8 29.8 0 0 1 43.18 128M240 120c-23 0-24-11.69-24-36c0-11.72 0-23.84-4.8-33.74C205.32 38.14 193.48 32 176 32a8 8 0 0 0 0 16c23 0 24 11.69 24 36c0 11.72 0 23.84 4.8 33.74a29.8 29.8 0 0 0 8 10.26a29.8 29.8 0 0 0-8 10.26c-4.8 9.9-4.8 22-4.8 33.74c0 24.31-1 36-24 36a8 8 0 0 0 0 16c17.48 0 29.32-6.14 35.2-18.26c4.8-9.9 4.8-22 4.8-33.74c0-24.31 1-36 24-36a8 8 0 0 0 0-16"/>
    </g>
    
    <g id="table-icon">
      <path fill="currentColor" d="M5 21q-.825 0-1.412-.587T3 19V5q0-.825.588-1.412T5 3h14q.825 0 1.413.588T21 5v14q0 .825-.587 1.413T19 21zm6-6H5v4h6zm2 0v4h6v-4zm-2-2V9H5v4zm2 0h6V9h-6zM5 7h14V5H5z"/>
    </g>
    
    <g id="local-storage-icon">
      <path fill="currentColor" d="M15,9H5V5H15M12,19A3,3 0 0,1 9,16A3,3 0 0,1 12,13A3,3 0 0,1 15,16A3,3 0 0,1 12,19M17,3H5C3.89,3 3,3.9 3,5V9A2,2 0 0,0 5,11V19A2,2 0 0,0 7,21H17A2,2 0 0,0 19,19V5C19,3.89 18.1,3 17,3Z"/>
    </g>
    
    <g id="session-storage-icon">
      <path fill="currentColor" d="M6 2v6h.01L6 8.01L10 12l-4 4-.01.01V22h12v-5.99l-.01-.01L14 12l4-3.99V8h-.01L18 2H6zm2 2h8v2.5L12 10 8 6.5V4zm8 16v-2.5L12 14l-4 3.5V20h8z"/>
    </g>
  </defs>
</svg>

<div id="collapsed">
  <h1>
    <svg viewBox="0 0 24 24" width="24" height="24" xmlns="http://www.w3.org/2000/svg">
      <use href="#datastar-logo"/>
    </svg>
  </h1>
</div>
<header style="display: none">
  <h1>
    <button aria-label="Collapse">\xBB</button>
    Datastar Inspector
    <svg viewBox="0 0 24 24" width="20" height="20" xmlns="http://www.w3.org/2000/svg">
      <use href="#datastar-logo"/>
    </svg>
  </h1>
  <nav>
    <button data-type="currentSignals" aria-selected="true" aria-label="Switch to current signals tab" title="Current signals with the ability to filter">
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"><path fill="currentColor" fill-rule="evenodd" d="M12.442 13.033c-.278.307-.319.777-.05 1.092c.27.314.747.353 1.033.053a7.5 7.5 0 1 0-10.85 0c.286.3.763.261 1.032-.053c.27-.315.23-.785-.05-1.092a6 6 0 1 1 8.884 0m-.987-1.15c-.265.318-.745.279-1.015-.036c-.27-.314-.223-.784.015-1.123a3 3 0 1 0-4.91 0c.238.339.284.809.015 1.123c-.27.315-.75.354-1.015.036a4.5 4.5 0 1 1 6.91 0M8 10.5a1.5 1.5 0 1 0 0-3a1.5 1.5 0 0 0 0 3" clip-rule="evenodd"/></svg>
      <span id="currentSignalsCount">0</span>
    </button>
    <button data-type="signalPatchEvent" aria-selected="false" aria-label="Switch to signal patch events tab" title="Signal patch events">
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M5 3h2v2H5v5a2 2 0 0 1-2 2a2 2 0 0 1 2 2v5h2v2H5c-1.07-.27-2-.9-2-2v-4a2 2 0 0 0-2-2H0v-2h1a2 2 0 0 0 2-2V5a2 2 0 0 1 2-2m14 0a2 2 0 0 1 2 2v4a2 2 0 0 0 2 2h1v2h-1a2 2 0 0 0-2 2v4a2 2 0 0 1-2 2h-2v-2h2v-5a2 2 0 0 1 2-2a2 2 0 0 1-2-2V5h-2V3zm-7 12a1 1 0 0 1 1 1a1 1 0 0 1-1 1a1 1 0 0 1-1-1a1 1 0 0 1 1-1m-4 0a1 1 0 0 1 1 1a1 1 0 0 1-1 1a1 1 0 0 1-1-1a1 1 0 0 1 1-1m8 0a1 1 0 0 1 1 1a1 1 0 0 1-1 1a1 1 0 0 1-1-1a1 1 0 0 1 1-1"/></svg>
      <span id="signalPatchEventCount">0</span>
    </button>
    <button data-type="sseEvent" aria-selected="false" aria-label="Switch to SSE tab" title="SSE events">
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="m12 17.56l4.07-1.13l.55-6.1H9.38L9.2 8.3h7.6l.2-1.99H7l.56 6.01h6.89l-.23 2.58l-2.22.6l-2.22-.6l-.14-1.66h-2l.29 3.19zM4.07 3h15.86L18.5 19.2L12 21l-6.5-1.8z"/></svg>
      <span id="sseEventCount">0</span>
    </button>
    <button data-type="persistedData" aria-selected="false" aria-label="Switch to persisted data tab" title="Persisted data from localStorage and sessionStorage">
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M4 2h14v2H4v16h2v-6h12v6h2V6h2v16H2V2zm4 18h8v-4H8zM20 6h-2V4h2zM6 6h9v4H6z"/></svg>
      <span id="persistedDataCount">0</span>
    </button>
  </nav>
</header>
<main style="display: none">
  <section id="currentSignalsContent" class="currentSignals">
    <ul>
      <li>
        <pre></pre>
      </li>
    </ul>
    <footer class="filters">
      <div class="filter-row">
        <div class="input-with-icon">
          <input type="text" id="currentSignalsIncludeFilter" placeholder="Include filter (default: .*)" title="Only show signals matching this pattern. Use regex or exact match (with wildcard support). Default '.*' matches all signals." />
        </div>
        <div class="input-with-icon">
          <input type="text" id="currentSignalsExcludeFilter" placeholder="Exclude filter (default: (^|\\.)_)" title="Hide signals matching this pattern. Use regex or exact match (with wildcard support). Default '(^|\\.)_' hides signals starting with underscore or containing '._'." />
        </div>
        <button id="resetFilters" title="Reset filters">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24">
              <use href="#reset-icon"/>
            </svg>
          </button>
        <label title="Enable wildcard matching (\`*\` matches any character)">
          <input type="checkbox" id="currentSignalsExactMatchCheckbox" />
          <span>Wildcard</span>
        </label>
      </div>
      <div class="filter-row-secondary">
        <div class="regex-display" id="currentSignalsRegexDisplay" style="display: none">
          <button id="copyFilterObject" title="Copy filter as regex object">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24">
              <use href="#copy-icon"/>
            </svg>
            <span style="display: none"></span>
          </button>
          <span class="regex-content">
            <span class="regex-include"></span>
            <span class="regex-exclude"></span>
          </span>
        </div>
        <label title="Toggle between object and path view">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 256 256">
            <use href="#tree-icon"/>
          </svg>
          <div class="toggle-switch">
            <input type="checkbox" id="currentSignalsTableViewCheckbox" />
            <span class="toggle-slider"></span>
          </div>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
            <use href="#table-icon"/>
          </svg>
        </label>
      </div>
    </footer>
  </section>
  <section id="signalPatchEventContent" class="signalPatchEvents" style="display: none">
    <p>No events yet.</p>
    <ul></ul>
    <footer>
      <label for="signalPatchMaxEvents">Max events:</label>
      <input type="number" id="signalPatchMaxEvents" placeholder="Max events" min="1" style="width: 80px;" />
      <button data-type="signalPatchEvent" style="display: none">Clear</button>
    </footer>
  </section>
  <section id="sseEventContent" class="sseEvents" style="display: none">
    <p>No events yet.</p>
    <ul></ul>
    <footer>
      <label for="sseMaxEvents">Max events:</label>
      <input type="number" id="sseMaxEvents" placeholder="Max events" min="1" style="width: 80px;" />
      <button data-type="sseEvent" style="display: none">Clear</button>
    </footer>
  </section>
  <section id="persistedDataContent" class="persistedData" style="display: none">
    <h3>
      <span id="persistedDataStorageTitle">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" style="vertical-align: middle;"><use href="#local-storage-icon"/></svg>
        Local Storage
      </span>
      <button id="refreshPersistedData" title="Refresh persisted data">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24">
          <use href="#refresh-icon"/>
        </svg>
        <span style="display: none"></span>
      </button>
    </h3>
    <p>No persisted data found.</p>
    <ul>
      <li>
        <pre></pre>
      </li>
    </ul>
    <footer class="filters">
      <div class="filter-row">
        <div class="input-with-icon">
          <input type="text" id="persistedDataIncludeFilter" placeholder="Include filter (default: .*)" title="Only show persisted data matching this pattern. Use regex or exact match (with wildcard support). Default '.*' matches all data." />
        </div>
        <div class="input-with-icon">
          <input type="text" id="persistedDataExcludeFilter" placeholder="Exclude filter (default: (^|\\.)_)" title="Hide persisted data matching this pattern. Use regex or exact match (with wildcard support). Default '(^|\\.)_' hides keys starting with underscore or containing '._'." />
        </div>
        <button id="resetPersistedDataFilters" title="Reset filters">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24">
            <use href="#reset-icon"/>
          </svg>
        </button>
        <label title="Enable wildcard matching (\`*\` matches any character)">
          <input type="checkbox" id="persistedDataExactMatchCheckbox" />
          <span>Wildcard</span>
        </label>
      </div>
      <div class="filter-row-secondary">
        <div class="regex-display" id="persistedDataRegexDisplay">
          <button id="copyPersistedDataFilterObject" title="Copy filter as regex object">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24">
              <use href="#copy-icon"/>
            </svg>
            <span style="display: none"></span>
          </button>
          <span class="regex-content">
            <span class="regex-include"></span>
            <span class="regex-exclude"></span>
          </span>
        </div>
        <label title="Toggle between localStorage and sessionStorage">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
            <use href="#local-storage-icon"/>
          </svg>
          <div class="toggle-switch">
            <input type="checkbox" id="persistedDataStorageToggle" />
            <span class="toggle-slider"></span>
          </div>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
            <use href="#session-storage-icon"/>
          </svg>
        </label>
        <label title="Toggle between object and path view">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 256 256">
            <use href="#tree-icon"/>
          </svg>
          <div class="toggle-switch">
            <input type="checkbox" id="persistedDataTableViewCheckbox" />
            <span class="toggle-slider"></span>
          </div>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
            <use href="#table-icon"/>
          </svg>
        </label>
      </div>
    </footer>
  </section>
</main>`;var z=class extends HTMLElement{constructor(){super();this.expanded=!1;this.open=!1;this.tab="currentSignals";this.maxEventsVisible=20;this.signalPatchEventCount=0;this.sseEventCount=0;this.signalPatchEvents=void 0;this.sseEvents=void 0;this.currentSignalsIncludeFilter="";this.currentSignalsExcludeFilter="";this.currentSignalsUseExactMatch=!1;this.currentSignalsTableView=!1;this.currentSignals={};this.persistedDataIncludeFilter="";this.persistedDataExcludeFilter="";this.persistedDataUseExactMatch=!1;this.persistedDataTableView=!1;this.persistedDataUseSessionStorage=!1;this.persistedData={};this.persistedDataCount=0;this.highlightedElements=[];this.highlightStyleElement=null;this.attachShadow({mode:"open"});let e=sessionStorage.getItem("datastar-inspector-state");if(e){let t=JSON.parse(e);this.expanded=t.expanded??this.expanded,this.open=t.open??this.open,this.tab=t.tab??this.tab,this.currentSignalsIncludeFilter=t.currentSignalsIncludeFilter??this.currentSignalsIncludeFilter,this.currentSignalsExcludeFilter=t.currentSignalsExcludeFilter??this.currentSignalsExcludeFilter,this.currentSignalsUseExactMatch=t.currentSignalsUseExactMatch??this.currentSignalsUseExactMatch,this.currentSignalsTableView=t.currentSignalsTableView??this.currentSignalsTableView,this.persistedDataIncludeFilter=t.persistedDataIncludeFilter??this.persistedDataIncludeFilter,this.persistedDataExcludeFilter=t.persistedDataExcludeFilter??this.persistedDataExcludeFilter,this.persistedDataUseExactMatch=t.persistedDataUseExactMatch??this.persistedDataUseExactMatch,this.persistedDataTableView=t.persistedDataTableView??this.persistedDataTableView,this.persistedDataUseSessionStorage=t.persistedDataUseSessionStorage??this.persistedDataUseSessionStorage}let s="el.setAttribute('data-current-signals', JSON.stringify($))";for(let t of["data-on-load","data-on-signal-patch"])this.setAttribute(t,s)}static get observedAttributes(){return["data-current-signals","max-events-visible"]}connectedCallback(){this.render(),this.injectHighlightStyles(),document.addEventListener("datastar-fetch",this.handleSseEvent.bind(this)),document.addEventListener("datastar-signal-patch",this.handleSignalPatchEvent.bind(this)),window.addEventListener("storage",this.handleStorageChange.bind(this));let e=this.getAttribute("max-events-visible");if(e){let s=Number.parseInt(e);!Number.isNaN(s)&&s>0&&(this.maxEventsVisible=s)}this.updatePersistedDataCount()}disconnectedCallback(){this.removeEventListener("datastar-fetch",this.handleSseEvent),this.removeEventListener("datastar-signal-patch",this.handleSignalPatchEvent),window.removeEventListener("storage",this.handleStorageChange),this.removeHighlightStyles()}applyTheme(){let e=window.matchMedia("(prefers-color-scheme: dark)").matches;this.shadowRoot&&(this.shadowRoot.host.classList.toggle("dark-theme",e),this.shadowRoot.host.setAttribute("data-theme",e?"dark":"light")),window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change",s=>{this.shadowRoot&&(this.shadowRoot.host.classList.toggle("dark-theme",s.matches),this.shadowRoot.host.setAttribute("data-theme",s.matches?"dark":"light"))})}attributeChangedCallback(e,s,t){e==="data-current-signals"&&s!==t&&t!==null&&this.handleCurrentSignalsChange(t),e==="max-events-visible"&&s!==t&&t!==null&&(this.maxEventsVisible=Number.parseInt(t),this.garbageCollectEvents(this.shadowRoot?.querySelector("#sseEventContent ul")),this.garbageCollectEvents(this.shadowRoot?.querySelector("#signalPatchEventContent ul")))}render(){if(!this.shadowRoot)return;this.shadowRoot.innerHTML=L,this.applyTheme(),this.renderExpand(),this.renderOpen(),this.renderTab();for(let t of this.shadowRoot.querySelectorAll("#collapsed, h1 button"))t.addEventListener("click",()=>this.toggleExpand());for(let t of this.shadowRoot.querySelectorAll("#collapsed, header"))t.addEventListener("click",()=>this.toggleOpen());for(let t of this.shadowRoot.querySelectorAll("header h1 button, header nav"))t.addEventListener("click",a=>a.stopPropagation());for(let t of this.shadowRoot.querySelectorAll("header nav button"))t.addEventListener("click",()=>this.selectTab(t.dataset.type||""));for(let t of this.shadowRoot.querySelectorAll("footer button[data-type]"))t.addEventListener("click",()=>this.clearEvents(t.dataset.type||""));let e=this.shadowRoot.querySelector("#signalPatchMaxEvents"),s=this.shadowRoot.querySelector("#sseMaxEvents");e&&(e.value=this.maxEventsVisible.toString(),e.addEventListener("input",()=>{let t=Number.parseInt(e.value);!Number.isNaN(t)&&t>0&&(this.maxEventsVisible=t,s&&(s.value=t.toString()),this.garbageCollectEvents(this.shadowRoot?.querySelector("#signalPatchEventContent ul")),this.garbageCollectEvents(this.shadowRoot?.querySelector("#sseEventContent ul")))})),s&&(s.value=this.maxEventsVisible.toString(),s.addEventListener("input",()=>{let t=Number.parseInt(s.value);!Number.isNaN(t)&&t>0&&(this.maxEventsVisible=t,e&&(e.value=t.toString()),this.garbageCollectEvents(this.shadowRoot?.querySelector("#signalPatchEventContent ul")),this.garbageCollectEvents(this.shadowRoot?.querySelector("#sseEventContent ul")))})),this.setupFilterHandlers(),this.setupPersistedDataHandlers(),this.tab==="persistedData"&&this.updatePersistedDataDisplay()}setupFilterHandlers(){if(!this.shadowRoot)return;let e=this.shadowRoot.querySelector("#currentSignalsIncludeFilter"),s=this.shadowRoot.querySelector("#currentSignalsExcludeFilter"),t=this.shadowRoot.querySelector("#currentSignalsExactMatchCheckbox"),a=this.shadowRoot.querySelector("#resetFilters"),i=this.shadowRoot.querySelector("#copyFilterObject"),n=this.shadowRoot.querySelector("#currentSignalsTableViewCheckbox");e&&(e.value=this.currentSignalsIncludeFilter||(this.currentSignalsUseExactMatch?"":".*"),e.addEventListener("input",()=>{this.currentSignalsIncludeFilter=e.value,this.saveState(),this.updateCurrentSignalsDisplay(),this.updateRegexDisplay()})),s&&(s.value=this.currentSignalsExcludeFilter||(this.currentSignalsUseExactMatch?"":"(^|\\.)_"),s.addEventListener("input",()=>{this.currentSignalsExcludeFilter=s.value,this.saveState(),this.updateCurrentSignalsDisplay(),this.updateRegexDisplay()})),t&&(t.checked=this.currentSignalsUseExactMatch,t.addEventListener("change",()=>{this.currentSignalsUseExactMatch=t.checked,this.saveState(),this.updateCurrentSignalsDisplay(),this.updateRegexDisplay()})),a&&a.addEventListener("click",()=>{this.currentSignalsIncludeFilter=".*",this.currentSignalsExcludeFilter="(^|\\.)_",this.currentSignalsUseExactMatch=!1,e.value=this.currentSignalsIncludeFilter,s.value=this.currentSignalsExcludeFilter,t.checked=this.currentSignalsUseExactMatch,this.saveState(),this.updateCurrentSignalsDisplay(),this.updateRegexDisplay()}),i&&i.addEventListener("click",()=>{this.copyFilterObject(i)}),n&&(n.checked=this.currentSignalsTableView,n.addEventListener("change",()=>{this.currentSignalsTableView=n.checked,this.saveState(),this.updateCurrentSignalsDisplay()})),this.updateRegexDisplay()}handleCurrentSignalsChange(e){this.currentSignals=JSON.parse(e),this.updateCurrentSignalsDisplay()}exactMatchAsRegex(e){return RegExp(`^${e.replace(/\./g,"\\.").replace(/\*/g,".*")}$`)}flattenObject(e,s=""){let t=[];for(let[a,i]of Object.entries(e)){let n=s?`${s}.${a}`:a;i!==null&&typeof i=="object"?Array.isArray(i)?(t.push([n,i]),i.forEach((r,o)=>{r!==null&&typeof r=="object"?t.push(...this.flattenObject(r,`${n}.${o}`)):t.push([`${n}.${o}`,r])})):(t.push([n,i]),t.push(...this.flattenObject(i,n))):t.push([n,i])}return t}updateRegexDisplay(){if(!this.shadowRoot)return;let e=this.shadowRoot.querySelector("#currentSignalsRegexDisplay"),s=e?.querySelector(".regex-include"),t=e?.querySelector(".regex-exclude");if(!e||!s||!t)return;e.style.display="";let a=this.currentSignalsIncludeFilter||".*",i=this.currentSignalsExcludeFilter||"(^|\\.)_";s.textContent=this.currentSignalsUseExactMatch&&this.currentSignalsIncludeFilter?this.exactMatchAsRegex(a).toString():`/${a}/`,t.textContent=this.currentSignalsUseExactMatch&&this.currentSignalsExcludeFilter?this.exactMatchAsRegex(i).toString():`/${i}/`}copyFilterObject(e){let s=this.currentSignalsIncludeFilter||".*",t=this.currentSignalsExcludeFilter||"(^|\\.)_",a={include:this.currentSignalsUseExactMatch?this.exactMatchAsRegex(s).toString():`/${s}/`,exclude:this.currentSignalsUseExactMatch?this.exactMatchAsRegex(t).toString():`/${t}/`},i=`{ include: ${a.include}, exclude: ${a.exclude} }`,n=e.querySelector("span");navigator.clipboard.writeText(i).then(()=>{n.textContent="Copied",n.style.display=""}).catch(()=>{n.textContent="Copy failed",n.style.display=""}).finally(()=>{setTimeout(()=>{n.style.display="none"},2e3)})}updateCurrentSignalsDisplay(){this.currentSignalsTableView?this.updateCurrentSignalsTableDisplay():this.updateCurrentSignalsTreeDisplay()}updateCurrentSignalsTableDisplay(){let e=this.currentSignalsIncludeFilter||this.currentSignalsExcludeFilter,s=this.flattenObject(this.currentSignals),t=0,n=`
      <table class="signals-table">
        <thead>
          <tr>
            <th>Path</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          ${s.map(([o,l])=>{let p=this.matchesFilter(o,this.currentSignalsIncludeFilter,this.currentSignalsExcludeFilter,this.currentSignalsUseExactMatch);return p&&t++,{path:o,value:l,matches:p}}).map(({path:o,value:l,matches:p})=>{let v=JSON.stringify(l),d=e?p?"matched":"filtered":"";return`<tr${d?` class="${d}"`:""}><td>${this.escapeHtml(o)}</td><td>${this.escapeHtml(v)}</td></tr>`}).join(`
`)}
        </tbody>
      </table>
    `;this.setTextContent("#currentSignalsCount",e?`${t}/${s.length}`:s.length);let r=this.shadowRoot?.querySelector("#currentSignalsContent ul li");r&&(r.innerHTML=`<div style="flex: 1; overflow-x: auto;">${n}</div>`,this.addCurrentSignalsButtons(r),this.addHoverHandlersToSignalPaths(r))}escapeHtml(e){let s=document.createElement("div");return s.textContent=e,s.innerHTML}updateCurrentSignalsTreeDisplay(){let e=this.currentSignalsIncludeFilter||this.currentSignalsExcludeFilter,s=new Set,t=0;if(e){let l=this.flattenObject(this.currentSignals);for(let[p]of l)if(this.matchesFilter(p,this.currentSignalsIncludeFilter,this.currentSignalsExcludeFilter,this.currentSignalsUseExactMatch)){t++;let v=p.split(".");for(let d=1;d<=v.length;d++)s.add(v.slice(0,d).join("."))}}let a=(l,p="",v="  ")=>{let d=[],g=Object.entries(l);return g.forEach(([x,m],u)=>{let c=p?`${p}.${x}`:x,h=u===g.length-1,f=s.has(c),b;if(e?b=f?'class="matched"':'class="filtered"':b=`data-path="${c}"`,m!==null&&typeof m=="object"&&!Array.isArray(m)){d.push(`${v}<span ${b}>"${x}": {`);let y=a(m,c,`${v}  `);d.push(...y),d.push(`${v}}${h?"</span>":",</span>"}`)}else{let D=JSON.stringify(m,null,2).split(`
`);d.push(`${v}<span ${b}>"${x}": ${this.escapeHtml(D[0])}`);for(let E=1;E<D.length;E++)d.push(`${v}${this.escapeHtml(D[E])}`);d[d.length-1]+=h?"</span>":",</span>"}}),d},i=[],n=e&&t>0;i.push(n?'<span class="matched">{</span>':"{"),i.push(...a(this.currentSignals)),i.push(n?'<span class="matched">}</span>':"}");let r=this.flattenObject(this.currentSignals).length;this.setTextContent("#currentSignalsCount",e?`${t}/${r}`:r);let o=this.shadowRoot?.querySelector("#currentSignalsContent ul li");o&&(o.innerHTML=`<div style="flex: 1; overflow-x: auto;"><pre>${i.join(`
`)}</pre></div>`,this.addCurrentSignalsButtons(o),this.addHoverHandlersToSignalPaths(o))}handleSignalPatchEvent(e){let s=e.detail,t=JSON.stringify(s);this.signalPatchEventCount++,this.setTextContent("#signalPatchEventCount",this.signalPatchEventCount),this.appendEvent(this.signalPatchEventCount,"#signalPatchEventContent",t,JSON.stringify(s,null,2),s,!0)}handleSseEvent(e){let s=e.detail;if(!s.type.startsWith("datastar-"))return;this.sseEventCount++,this.setTextContent("#sseEventCount",this.sseEventCount);let t="",a=s.argsRaw;for(let[i,n]of Object.entries(a))t+=`data: ${i} ${n}
`;t=t.trimEnd(),this.appendEvent(this.sseEventCount,"#sseEventContent",`event: ${s.type}`,t,s,!1)}toggleExpand(){this.expanded=!this.expanded,this.renderExpand(),this.saveState()}renderExpand(){this.setVisibility("#collapsed",!this.expanded),this.setVisibility("header",this.expanded),this.expanded||(this.open=!0,this.toggleOpen())}toggleOpen(){this.open=!this.open,this.renderOpen(),this.saveState()}renderOpen(){this.setVisibility("main",this.open)}selectTab(e){this.open=!0,this.tab=e,this.renderTab(),this.saveState(),e==="persistedData"&&this.updatePersistedDataDisplay()}renderTab(){this.setAttributeValue("header button","aria-selected","false"),this.setAttributeValue(`header button[data-type=${this.tab}]`,"aria-selected","true"),this.setVisibility("main",this.open),this.setVisibility("main section",!1),this.setVisibility(`#${this.tab}Content`,!0),this.saveState()}saveState(){let e={expanded:this.expanded,open:this.open,tab:this.tab,currentSignalsIncludeFilter:this.currentSignalsIncludeFilter,currentSignalsExcludeFilter:this.currentSignalsExcludeFilter,currentSignalsUseExactMatch:this.currentSignalsUseExactMatch,currentSignalsTableView:this.currentSignalsTableView,persistedDataIncludeFilter:this.persistedDataIncludeFilter,persistedDataExcludeFilter:this.persistedDataExcludeFilter,persistedDataUseExactMatch:this.persistedDataUseExactMatch,persistedDataTableView:this.persistedDataTableView,persistedDataUseSessionStorage:this.persistedDataUseSessionStorage};sessionStorage.setItem("datastar-inspector-state",JSON.stringify(e))}clearEvents(e){e==="sse"?this.sseEventCount=0:this.signalPatchEventCount=0,this.setTextContent(`#${e}Count`,0),this.setInnerHtml(`#${e}Content ul`,""),this.setVisibility(`#${e}Content footer button[data-type]`,!1),this.setVisibility(`#${e}Content p`,!0)}setTextContent(e,s){let t=this.shadowRoot?.querySelector(e);t&&(t.textContent=`${s}`)}setInnerHtml(e,s){let t=this.shadowRoot?.querySelector(e);t&&(t.innerHTML=s)}setAttributeValue(e,s,t){for(let a of this.shadowRoot?.querySelectorAll(e)||[])a.setAttribute(s,t)}setVisibility(e,s){this.setAttributeValue(e,"style",s?"":"display: none;")}appendEvent(e,s,t,a,i,n){let r=this.shadowRoot?.querySelector(`${s} ul`);if(!r)return;let o=document.createElement("pre");o.textContent=a;let l=document.createElement("div");l.innerHTML=`
      <span>${e}.</span>
      <details>
        <summary class="${n?"blurrable":""}">${t}</summary>
        <pre>${o.innerHTML}</pre>
      </details>
    `;let d=(this.shadowRoot?.querySelector("#copy-button-template")).content.cloneNode(!0).querySelector("button"),g=d.querySelector("span");d.onclick=()=>{let f=t.startsWith("event:")?`${t}
${a}`:a;navigator.clipboard.writeText(f).then(()=>{g.textContent="Copied",g.style.display=""}).catch(()=>{g.textContent="Copy Failed",g.style.display=""}).finally(()=>{setTimeout(()=>{g.style.display="none"},2e3)})};let u=(this.shadowRoot?.querySelector("#log-button-template")).content.cloneNode(!0).querySelector("button"),c=u.querySelector("span");u.onclick=()=>{console.log(i),c.textContent="Logged",c.style.display="",setTimeout(()=>{c.style.display="none"},2e3)};let h=document.createElement("li");h.appendChild(l),h.appendChild(d),h.appendChild(u),r.prepend(h),this.garbageCollectEvents(r),this.setVisibility(`${s} footer button[data-type]`,!0),this.setVisibility(`${s} p`,!1)}createButton(e,s){let t=document.createElement("button");return t.textContent=e,t.onclick=s,t}garbageCollectEvents(e){if(e)for(;e.children.length>this.maxEventsVisible;)e.removeChild(e.lastChild)}addCurrentSignalsButtons(e){let a=(this.shadowRoot?.querySelector("#copy-button-template")).content.cloneNode(!0).querySelector("button"),i=a.querySelector("span");a.onclick=()=>{let p=JSON.stringify(this.currentSignals,null,2);navigator.clipboard.writeText(p).then(()=>{i.textContent="Copied",i.style.display=""}).catch(()=>{i.textContent="Copy Failed",i.style.display=""}).finally(()=>{setTimeout(()=>{i.style.display="none"},2e3)})};let o=(this.shadowRoot?.querySelector("#log-button-template")).content.cloneNode(!0).querySelector("button"),l=o.querySelector("span");o.onclick=()=>{console.log(this.currentSignals),l.textContent="Logged",l.style.display="",setTimeout(()=>{l.style.display="none"},2e3)},e.appendChild(a),e.appendChild(o),this.addHoverHandlersToSignalPaths(e)}matchesFilter(e,s,t,a){let i=s||".*",n=t||"(^|\\.)_";if(a){let r=this.exactMatchAsRegex(i),o=this.exactMatchAsRegex(n);return r.test(e)&&!o.test(e)}try{let r=new RegExp(i),o=new RegExp(n);return r.test(e)&&!o.test(e)}catch{return e.includes(i)&&!e.includes(n)}}injectHighlightStyles(){this.highlightStyleElement||(this.highlightStyleElement=document.createElement("style"),this.highlightStyleElement.textContent=`
      @keyframes datastar-highlight-pulse {
        0% {
          box-shadow: 
            0 0 0 2px rgba(165, 153, 255, 0.8),
            0 0 0 6px rgba(165, 153, 255, 0.4),
            0 0 16px 4px rgba(165, 153, 255, 0.6);
          background-color: rgba(165, 153, 255, 0.2);
        }
        50% {
          box-shadow: 
            0 0 0 4px rgba(165, 153, 255, 1),
            0 0 0 12px rgba(165, 153, 255, 0.6),
            0 0 32px 8px rgba(165, 153, 255, 0.8);
          background-color: rgba(165, 153, 255, 0.4);
        }
        100% {
          box-shadow: 
            0 0 0 2px rgba(165, 153, 255, 0.8),
            0 0 0 6px rgba(165, 153, 255, 0.4),
            0 0 16px 4px rgba(165, 153, 255, 0.6);
          background-color: rgba(165, 153, 255, 0.2);
        }
      }
      
      .datastar-signal-highlight {
        animation: datastar-highlight-pulse 0.8s ease-in-out infinite !important;
        position: relative !important;
        z-index: 9998 !important;
      }
    `,document.head.appendChild(this.highlightStyleElement))}removeHighlightStyles(){this.highlightStyleElement&&(document.head.removeChild(this.highlightStyleElement),this.highlightStyleElement=null)}findElementsUsingSignal(e){let s=[],t=this.generateSignalPatterns(e),a=document.querySelectorAll("*");for(let i of a){let n=i.attributes;for(let r of n)if(r.name.startsWith("data-")&&this.attributeUsesSignal(r.value,t)){s.push(i);break}}return s}generateSignalPatterns(e){let s=[],t=e.split(".");for(let a=t.length;a>0;a--){let i=t.slice(0,a).join(".");s.push(`$${i}`)}return s}attributeUsesSignal(e,s){for(let t of s)if(new RegExp(`\\${t.replace(".","\\.")}(?![a-zA-Z0-9_])`,"g").test(e))return!0;return!1}highlightElements(e){this.clearHighlights();for(let s of e)s.classList.add("datastar-signal-highlight"),this.highlightedElements.push(s)}clearHighlights(){for(let e of this.highlightedElements)e.classList.remove("datastar-signal-highlight");this.highlightedElements=[]}addHoverHandlersToSignalPaths(e,s="current"){let t=this.extractSignalPaths(e,s);for(let a of t)a.element.classList.add("signal-path-hoverable"),a.element.addEventListener("mouseenter",()=>{let i=this.findElementsUsingSignal(a.path);this.highlightElements(i)}),a.element.addEventListener("mouseleave",()=>{this.clearHighlights()})}extractSignalPaths(e,s){let t=[],a=s==="current"?this.currentSignalsTableView:this.persistedDataTableView,i=s==="current"?this.currentSignals:this.persistedData;if(a){let n=e.querySelectorAll("td:first-child");for(let r of n){let o=r.textContent?.trim();o&&t.push({element:r,path:o})}}else{let n=s==="persisted"?Array.from(e.querySelectorAll("pre")):[e.querySelector("pre")].filter(r=>r!==null);for(let r of n){let o=r.querySelectorAll("span"),l=new Set;for(let d of o){let g=d.getAttribute("data-path");if(g&&!l.has(g)){l.add(g),t.push({element:d,path:g});continue}let m=(d.textContent||"").match(/"([^"]+)":\s*/);if(m){let u=m[1];if(s==="current"){let c=this.flattenObject(i);for(let[h]of c)if(h.endsWith(u)||h===u){l.has(h)||(l.add(h),t.push({element:d,path:h}));break}}else for(let[,c]of Object.entries(i)){let h=this.flattenObject(c);for(let[f]of h)if(f.endsWith(u)||f===u){l.has(f)||(l.add(f),t.push({element:d,path:f}));break}}}}let p=document.createTreeWalker(r,NodeFilter.SHOW_TEXT,{acceptNode:d=>d.parentElement?.tagName==="SPAN"?NodeFilter.FILTER_SKIP:NodeFilter.FILTER_ACCEPT}),v;for(;v=p.nextNode();){let g=(v.textContent||"").match(/"([^"]+)":\s*/);if(g){let x=g[1],m=v.parentElement;if(m){let u=new Set;if(s==="current"){let c=this.flattenObject(i);for(let[h]of c)if(h.endsWith(x)||h===x){u.has(h)||(u.add(h),t.push({element:m,path:h}));break}}else for(let[,c]of Object.entries(i)){let h=this.flattenObject(c);for(let[f]of h)if(f.endsWith(x)||f===x){u.has(f)||(u.add(f),t.push({element:m,path:f}));break}}}}}}}return t}handleStorageChange(){this.updatePersistedDataCount(),this.tab==="persistedData"&&this.updatePersistedDataDisplay()}updatePersistedDataCount(){let e=this.persistedDataUseSessionStorage?sessionStorage:localStorage,t=this.getDatastarPersistKeys().filter(a=>e.getItem(a)!==null).length;this.setTextContent("#persistedDataCount",t)}getDatastarPersistKeys(){let e=new Set,s=document.querySelectorAll("*");for(let t of s)for(let a of t.attributes)if(a.name.startsWith("data-persist")){let i=a.name;if(i==="data-persist")e.add("datastar");else if(i.startsWith("data-persist-")){let n=i.substring(13),[r]=n.split("__");r&&e.add(r)}}return Array.from(e)}updatePersistedDataDisplay(){let e=this.persistedDataUseSessionStorage?sessionStorage:localStorage,s=this.getDatastarPersistKeys(),t={},a=[];for(let n of s){let r=e.getItem(n);if(r!==null){a.push(n);try{t[n]=JSON.parse(r)}catch{t[n]=r}}}this.persistedData=t,this.persistedDataCount=a.length,this.setTextContent("#persistedDataCount",this.persistedDataCount);let i=this.shadowRoot?.querySelector("#persistedDataContent p");if(i&&(i.style.display=a.length===0?"":"none"),a.length===0){let n=this.shadowRoot?.querySelector("#persistedDataContent ul");n&&(n.innerHTML="");return}this.persistedDataTableView?this.updatePersistedDataTableDisplay():this.updatePersistedDataTreeDisplay()}updatePersistedDataTableDisplay(){let e=this.persistedDataUseSessionStorage?sessionStorage:localStorage,s=this.getDatastarPersistKeys(),t={};for(let l of s){let p=e.getItem(l);if(p!==null)try{t[l]=JSON.parse(p)}catch{t[l]=p}}let a=this.persistedDataIncludeFilter||this.persistedDataExcludeFilter,i=0,n=0,r=Object.entries(t).map(([l,p])=>{let v=this.flattenObject(p);n+=v.length;let g=v.map(([x,m])=>{let u=this.matchesFilter(x,this.persistedDataIncludeFilter,this.persistedDataExcludeFilter,this.persistedDataUseExactMatch);return u&&i++,{path:x,value:m,matches:u}}).map(({path:x,value:m,matches:u})=>{let c=JSON.stringify(m),h=a?u?"matched":"filtered":"";return`<tr${h?` class="${h}"`:""}><td>${this.escapeHtml(x)}</td><td>${this.escapeHtml(c)}</td></tr>`}).join(`
`);return`
          <div class="persisted-key-section">
            <div class="persisted-key-caption">
              <div class="caption-text">${this.escapeHtml(l)}</div>
              <div class="caption-buttons" data-key="${this.escapeHtml(l)}"></div>
            </div>
            <table>
              <thead>
                <tr>
                  <th class="thin">Path</th>
                  <th>Value</th>
                </tr>
              </thead>
              <tbody>
                ${g}
              </tbody>
            </table>
          </div>
        `}).join(`
`);this.setTextContent("#persistedDataCount",a?`${i}/${n}`:n);let o=this.shadowRoot?.querySelector("#persistedDataContent ul");if(o){o.innerHTML=`<li><div style="flex: 1; overflow-x: auto;">${r}</div></li>`;let l=o.querySelector("li");l&&this.addPersistedDataButtons(l)}}updatePersistedDataTreeDisplay(){let e=this.persistedDataUseSessionStorage?sessionStorage:localStorage,s=this.getDatastarPersistKeys(),t={};for(let l of s){let p=e.getItem(l);if(p!==null)try{t[l]=JSON.parse(p)}catch{t[l]=p}}let a=this.persistedDataIncludeFilter||this.persistedDataExcludeFilter,i=0,n=0,r=Object.entries(t).map(([l,p])=>{let v=this.flattenObject(p);n+=v.length;let d=new Set;if(a){for(let[m]of v)if(this.matchesFilter(m,this.persistedDataIncludeFilter,this.persistedDataExcludeFilter,this.persistedDataUseExactMatch)){i++;let u=m.split(".");for(let c=1;c<=u.length;c++)d.add(u.slice(0,c).join("."))}}let g=(m,u="",c="  ")=>{let h=[],f=Object.entries(m);return f.forEach(([b,y],D)=>{let E=u?`${u}.${b}`:b,C=D===f.length-1,T=d.has(E);if(y!==null&&typeof y=="object"&&!Array.isArray(y)){a?T?h.push(`${c}<span class="matched">"${b}": {`):h.push(`${c}<span class="filtered">"${b}": {`):h.push(`${c}"${b}": {`);let M=g(y,E,`${c}  `);h.push(...M),a?h.push(`${c}}${C?"</span>":",</span>"}`):h.push(`${c}}${C?"":","}`)}else{let S=JSON.stringify(y,null,2).split(`
`);if(a){T?h.push(`${c}<span class="matched">"${b}": ${this.escapeHtml(S[0])}`):h.push(`${c}<span class="filtered">"${b}": ${this.escapeHtml(S[0])}`);for(let w=1;w<S.length;w++)h.push(`${c}${S[w]}`);h[h.length-1]+=C?"</span>":",</span>"}else{h.push(`${c}"${b}": ${this.escapeHtml(S[0])}`);for(let w=1;w<S.length;w++)h.push(`${c}${this.escapeHtml(S[w])}`);C||(h[h.length-1]+=",")}}}),h},x=[];return x.push("{"),x.push(...g(p)),x.push("}"),`
          <div class="persisted-key-section">
            <div class="persisted-key-caption">
              <div class="caption-text">${this.escapeHtml(l)}</div>
              <div class="caption-buttons" data-key="${this.escapeHtml(l)}"></div>
            </div>
            <pre>${x.join(`
`)}</pre>
          </div>
        `}).join(`
`);this.setTextContent("#persistedDataCount",a?`${i}/${n}`:n);let o=this.shadowRoot?.querySelector("#persistedDataContent ul");if(o){o.innerHTML=`<li><div style="flex: 1; overflow-x: auto;">${r}</div></li>`;let l=o.querySelector("li");l&&this.addPersistedDataButtons(l)}}setupPersistedDataHandlers(){if(!this.shadowRoot)return;let e=this.shadowRoot.querySelector("#persistedDataIncludeFilter");e&&(this.persistedDataIncludeFilter||(this.persistedDataIncludeFilter=".*"),e.value=this.persistedDataIncludeFilter,e.addEventListener("input",()=>{this.persistedDataIncludeFilter=e.value,this.saveState(),this.updatePersistedDataDisplay(),this.updatePersistedDataRegexDisplay()}));let s=this.shadowRoot.querySelector("#persistedDataExcludeFilter");s&&(this.persistedDataExcludeFilter||(this.persistedDataExcludeFilter="(^|\\.)_"),s.value=this.persistedDataExcludeFilter,s.addEventListener("input",()=>{this.persistedDataExcludeFilter=s.value,this.saveState(),this.updatePersistedDataDisplay(),this.updatePersistedDataRegexDisplay()}));let t=this.shadowRoot.querySelector("#persistedDataExactMatchCheckbox");t&&(t.checked=this.persistedDataUseExactMatch,t.addEventListener("change",()=>{this.persistedDataUseExactMatch=t.checked,this.saveState(),this.updatePersistedDataDisplay(),this.updatePersistedDataRegexDisplay()}));let a=this.shadowRoot.querySelector("#persistedDataTableViewCheckbox");a&&(a.checked=this.persistedDataTableView,a.addEventListener("change",()=>{this.persistedDataTableView=a.checked,this.saveState(),this.updatePersistedDataDisplay()}));let i=this.shadowRoot.querySelector("#persistedDataStorageToggle");i&&(i.checked=this.persistedDataUseSessionStorage,i.addEventListener("change",()=>{this.persistedDataUseSessionStorage=i.checked,this.updateStorageTitle(),this.saveState(),this.updatePersistedDataCount(),this.updatePersistedDataDisplay()}));let n=this.shadowRoot.querySelector("#refreshPersistedData");n&&n.addEventListener("click",()=>{this.updatePersistedDataCount(),this.updatePersistedDataDisplay()});let r=this.shadowRoot.querySelector("#resetPersistedDataFilters");r&&r.addEventListener("click",()=>{this.persistedDataIncludeFilter=".*",this.persistedDataExcludeFilter="(^|\\.)_",this.persistedDataUseExactMatch=!1,e.value=this.persistedDataIncludeFilter,s.value=this.persistedDataExcludeFilter,t.checked=this.persistedDataUseExactMatch,this.saveState(),this.updatePersistedDataDisplay(),this.updatePersistedDataRegexDisplay()});let o=this.shadowRoot.querySelector("#copyPersistedDataFilterObject");o&&o.addEventListener("click",()=>{this.copyPersistedDataFilterObject(o)}),this.updatePersistedDataRegexDisplay(),this.updateStorageTitle()}updateStorageTitle(){if(!this.shadowRoot)return;let e=this.shadowRoot.querySelector("#persistedDataStorageTitle");e&&(this.persistedDataUseSessionStorage?e.innerHTML='<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" style="vertical-align: middle; margin-right: 4px;"><use href="#session-storage-icon"/></svg>Session Storage':e.innerHTML='<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" style="vertical-align: middle; margin-right: 4px;"><use href="#local-storage-icon"/></svg>Local Storage')}updatePersistedDataRegexDisplay(){if(!this.shadowRoot)return;let e=this.shadowRoot.querySelector("#persistedDataRegexDisplay"),s=e?.querySelector(".regex-include"),t=e?.querySelector(".regex-exclude");if(!e||!s||!t)return;let a=this.persistedDataIncludeFilter||".*",i=this.persistedDataExcludeFilter||"(^|\\.)_";s.textContent=this.persistedDataUseExactMatch&&this.persistedDataIncludeFilter?this.exactMatchAsRegex(a).toString():`/${a}/`,t.textContent=this.persistedDataUseExactMatch&&this.persistedDataExcludeFilter?this.exactMatchAsRegex(i).toString():`/${i}/`}copyPersistedDataFilterObject(e){let s=this.persistedDataIncludeFilter||".*",t=this.persistedDataExcludeFilter||"(^|\\.)_",a={include:this.persistedDataUseExactMatch?this.exactMatchAsRegex(s).toString():`/${s}/`,exclude:this.persistedDataUseExactMatch?this.exactMatchAsRegex(t).toString():`/${t}/`},i=`{ include: ${a.include}, exclude: ${a.exclude} }`,n=e.querySelector("span");navigator.clipboard.writeText(i).then(()=>{n.textContent="Copied",n.style.display=""}).catch(()=>{n.textContent="Copy failed",n.style.display=""}).finally(()=>{setTimeout(()=>{n.style.display="none"},2e3)})}addPersistedDataButtons(e){e.querySelectorAll(".caption-buttons").forEach(t=>{let a=t.getAttribute("data-key");if(!a||!this.persistedData[a])return;let i=this.persistedData[a],o=(this.shadowRoot?.querySelector("#copy-button-template")).content.cloneNode(!0).querySelector("button"),l=o.querySelector("span");o.onclick=()=>{let c=JSON.stringify(i,null,2);navigator.clipboard.writeText(c).then(()=>{l.textContent="Copied",l.style.display=""}).catch(()=>{l.textContent="Copy Failed",l.style.display=""}).finally(()=>{setTimeout(()=>{l.style.display="none"},2e3)})};let v=(this.shadowRoot?.querySelector("#log-button-template")).content.cloneNode(!0),d=v.querySelector("button"),g=d.querySelector("span");d.onclick=()=>{console.log(i),g.textContent="Logged",g.style.display="",setTimeout(()=>{g.style.display="none"},2e3)};let u=(this.shadowRoot?.querySelector("#clear-button-template")).content.cloneNode(!0).querySelector("button");u.onclick=()=>{confirm("Clear persisted data? This cannot be undone.")&&((this.persistedDataUseSessionStorage?sessionStorage:localStorage).removeItem(a),this.updatePersistedDataCount(),this.updatePersistedDataDisplay())},t.appendChild(o),t.appendChild(v),t.appendChild(u)}),this.addHoverHandlersToSignalPaths(e,"persisted")}};customElements.define("datastar-inspector",z);
/*! Datastar Inspector v1.0.3 */
//# sourceMappingURL=datastar-inspector.js.map
