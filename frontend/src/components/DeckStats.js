import React, { useMemo } from 'react';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend } from 'recharts';
import './DeckStats.css';

const COLORS = {
  'W': '#f8e7b9',
  'U': '#b3ceea',
  'B': '#a69f9d',
  'R': '#eb9f82',
  'G': '#c4d3ca',
  'Colorless': '#DCDCDC'
};

function DeckStats({ cards }) {
  // useMemo will only re-calculate the stats when the 'cards' prop changes.
  const { manaCurve, colorDistribution } = useMemo(() => {
    const curve = Array(8).fill(0).map((_, i) => ({ cmc: i, count: 0 })); // Bins for CMC 0 through 7+
    const colors = { 'W': 0, 'U': 0, 'B': 0, 'R': 0, 'G': 0, 'Colorless': 0 };

    cards.forEach(card => {
      // Mana Curve Calculation
      const cmc = Math.floor(card.cmc);
      if (cmc >= 7) {
        curve[7].count += card.quantity; // Group all 7+ CMC cards together
      } else if (cmc >= 0) {
        curve[cmc].count += card.quantity;
      }

      // Color Distribution Calculation
      if (card.colors && card.colors.length > 0) {
        card.colors.forEach(color => {
          if (colors.hasOwnProperty(color)) {
            colors[color] += card.quantity;
          }
        });
      } else {
        colors['Colorless'] += card.quantity;
      }
    });
    
    // Format for Pie Chart, filtering out empty color slices
    const pieData = Object.entries(colors)
      .filter(([_, value]) => value > 0)
      .map(([name, value]) => ({ name, value }));

    // Add '+' to the last label of the mana curve
    curve[7].cmc = '7+';

    return { manaCurve: curve, colorDistribution: pieData };
  }, [cards]);

  return (
    <div className="deck-stats-container">
      <div className="stat-chart">
        <h3>Mana Curve</h3>
        <ResponsiveContainer width="100%" height={250}>
          <BarChart data={manaCurve} margin={{ top: 5, right: 20, left: -10, bottom: 5 }}>
            <XAxis dataKey="cmc" stroke="#c0c0c0" />
            <YAxis allowDecimals={false} stroke="#c0c0c0" />
            <Tooltip
              contentStyle={{ backgroundColor: '#242424', border: '1px solid #444' }}
              cursor={{ fill: 'rgba(106, 13, 173, 0.2)' }}
            />
            <Bar dataKey="count" fill="#6a0dad" />
          </BarChart>
        </ResponsiveContainer>
      </div>
      <div className="stat-chart">
        <h3>Color Distribution</h3>
        <ResponsiveContainer width="100%" height={250}>
          <PieChart>
            <Pie
              data={colorDistribution}
              cx="50%"
              cy="50%"
              labelLine={false}
              outerRadius={80}
              fill="#8884d8"
              dataKey="value"
              nameKey="name"
              label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
            >
              {colorDistribution.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[entry.name]} />
              ))}
            </Pie>
            <Tooltip contentStyle={{ backgroundColor: '#242424', border: '1px solid #444' }} />
            <Legend />
          </PieChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}

export default DeckStats;
