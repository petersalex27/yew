# Yew Programming Language

Welcome to the Yew language repo!

Official Yew language site (not yet up as of 08/21/2024):
- yew-lang.org

## Installing
TODO

### Building from Source
TODO

## Commands
TODO

### Command `yew`: Yew Compiler
<table>
  <tr>
    <th>Option(s)</th>
    <th>Description</th>
    <th>Example</th>
    <th>Also</th>
  </tr>

  <tr>
    <th colspan="4"><h3><code>yew</code></h3></th>
  </tr>
  <tr>
    <td></td>
    <td>Starts repl interface</td>
    <td><code>yew</code></td>
    <td></td>
  </tr>

  <tr>
    <th colspan="4"><h3><code>yew repl</code></h3></th>
  </tr>
  <tr>
    <td></td>
    <td>Starts repl interface</td>
    <td><code>yew repl</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>-i [pkg1,pkg2,..]</code></td>
    <td>Imports pkg1, pkg2, ...</td>
    <td><code>yew repl -i base,reflect</code></td>
    <td><code>--import</code></td>
  </tr>
  <tr>
    <td><code>-L</code></td>
    <td>Runs in literate mode</td>
    <td><code>yew repl -L</code></td>
    <td><code>--lit, --literate</code></td>
  </tr>
  <tr>
    <td><code>-o &lt;file&gt;</code></td>
    <td>Outputs REPL input to <code>file</code></td>
    <td><code>yew repl -o record.yew</code></td>
    <td><code>--out, --output</code></td>
  </tr>

  <tr>
    <th colspan="4"><h3><code>yew build</code></h3></th>
  </tr>
  <tr>
    <td></td>
    <td>Builds package in pwd</td>
    <td><code>yew build</code></code>
    <td></td>
  </tr>
  <tr>
    <td><code>[pkg]</code></td>
    <td>Builds <code>pkg</code>. Must be first arg.</td>
    <td><code>yew build pkg</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>-o &lt;file&gt;</code></td>
    <td>Writes executable to <code>file</code></td>
    <td><code>yew build pkg -o a.out</code></td>
    <td><code>--out, --output</code></td>
  </tr>
  <tr>
    <td><code>-- &lt;pkg&gt;</code></td>
    <td>Builds package <code>pkg</code></td>
    <td><code>yew build -o a.out -- pkg</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>-i</code></td>
    <td>Stops after producing all IR</td>
    <td><code>yew build pkg -i</code></td>
    <td><code>--ir, --intermediate</td>
  </tr>
  <tr>
    <td><code>-w (all|none)</code></td>
    <td>Enables and disables all warnings resp.</td>
    <td><code>yew build pkg -w all</code>
    <td><code>--warning</code></td>
  </tr>
  <tr>
    <td><code>-w &lt;config&gt;</code></td>
    <td>Uses warning flags described in <code>config</code></td>
    <td><code>yew build pkg -w warn.config</code></td>
    <td><code>--warning</code></td>
  </tr>

  <tr>
    <th colspan="4"><h3><code>yew help</code></h3></th>
  </tr>
  <tr>
    <td></td>
    <td>Displays info for common commands</td>
    <td><code>yew help</code></td>
    <td><code>--common</code></td>
  </tr>
  <tr>
    <td><code>[topic]</code></td>
    <td>Displays info for <code>topic</code>. Must be first arg.</td>
    <td><code>yew help build</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>-b</code></td>
    <td>Displays info for builtins</td>
    <td><code>yew help Type -b</code></td>
    <td><code>--builtin, --builtins</code></td>
  </tr>
  <tr>
    <td><code>-o &lt;option&gt;</code></td>
    <td>Displays info for <code>option</code> of a command</td>
    <td><code>yew help build -o ir</code></td>
    <td><code>--opt, --option</td>
  </tr>
  <tr>
    <td><code>-- &lt;topic&gt;</td>
    <td>Displays info for <code>topic</code></td>
    <td><code>yew help -o ir -- build</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>-v [bool]</code></td>
    <td>Sets verbose help to <code>bool</code>. Default is <code>true</code>
    <td><code>yew help build -v true</code></td>
    <td><code>--verbose</code></td>
  </tr>

  <tr>
    <th colspan="4"><h3><code>yew version</code></h3></th>
  </tr>
  <tr>
    <td></td>
    <td>Displays running version of yew compiler</td>
    <td><code>yew version</code></td>
    <td></td>
  </tr>
</table>

TODO: finish

### Command: `ypk`: Yew Package Manager 
TODO

## License
Yew is distributed under the terms of the MIT license.

See the file `LICENSE` located in the same directory as this file for more details